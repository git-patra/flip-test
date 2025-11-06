package httpiface_test

import (
	"boilerplate-go/internal/pkg/statements/domain/constant"
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"
	repo "boilerplate-go/internal/pkg/statements/infrastructure/repo"
	httpiface "boilerplate-go/internal/pkg/statements/interfaces/http"
	"boilerplate-go/internal/pkg/statements/usecase"
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func buildRouter(r repo.InMemoryRepo, x bus.Exchange) http.Handler {
	parseCSV := usecase.NewParseCSVUsecase(r, x)
	getBalance := usecase.NewGetBalanceUsecase(r)
	getIssues := usecase.NewGetIssuesUsecase(r)
	h := httpiface.NewHandler(parseCSV, getBalance, getIssues)

	mux := chi.NewRouter()
	mux.Post("/statements", h.Upload)
	mux.Get("/balance", h.GetBalance)
	mux.Get("/transactions/issues", h.GetIssues)
	return mux
}

func multipartCSV(t *testing.T, csv string) (*bytes.Buffer, string) {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("file", "sample.csv")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write([]byte(csv)); err != nil {
		t.Fatalf("write csv: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	return &body, w.FormDataContentType()
}

// Test end-to-end: upload → parse(streaming) → balance & issues query
func TestUploadProcessQuery_Integration(t *testing.T) {
	// infra
	x := bus.NewExchange()
	r := repo.NewInMemoryRepo()

	// (opsional) start real consumers — tidak diperlukan untuk query,
	// tapi aman kalau ingin memastikan tidak error saat event dipublish.
	// Kita biarkan tidak ada consumer di test ini.

	srv := httptest.NewServer(buildRouter(r, x))
	defer srv.Close()

	csv := `timestamp,counterparty,type,amount,status,description
1674507883, JOHN DOE, DEBIT, 250000, SUCCESS, restaurant
1674507890, ACME CORP, CREDIT, 500000, SUCCESS, salary
1674507900, JANE, DEBIT, 125000, FAILED, atm
1674507910, SHOP, DEBIT, 85000, PENDING, groceries
`
	body, ctype := multipartCSV(t, csv)

	// POST /statements
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/statements", body)
	req.Header.Set("Content-Type", ctype)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("upload request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("upload status=%d", resp.StatusCode)
	}
	var up struct {
		Data struct {
			UploadID string `json:"upload_id"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&up); err != nil {
		t.Fatalf("decode upload resp: %v", err)
	}

	// GET /balance
	resp2, err := http.Get(srv.URL + "/balance?upload_id=" + up.Data.UploadID)
	if err != nil {
		t.Fatalf("balance request error: %v", err)
	}
	defer resp2.Body.Close()
	var bal struct {
		Data int64 `json:"data"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&bal); err != nil {
		t.Fatalf("decode balance: %v", err)
	}
	if bal.Data != 250000 {
		t.Fatalf("expected balance=250000 got %d", bal.Data)
	}

	// GET /transactions/issues
	resp3, err := http.Get(srv.URL + "/transactions/issues?upload_id=" + up.Data.UploadID + "&page=1&size=10")
	if err != nil {
		t.Fatalf("issues request error: %v", err)
	}
	defer resp3.Body.Close()
	var issues struct {
		Data     []entity.Transaction `json:"data"`
		Metadata struct {
			TotalPages int `json:"total_pages"`
			Page       int `json:"page"`
			Size       int `json:"size"`
			TotalData  int `json:"total_data"`
		}
	}
	if err := json.NewDecoder(resp3.Body).Decode(&issues); err != nil {
		t.Fatalf("decode issues: %v", err)
	}
	if issues.Metadata.TotalData != 2 {
		t.Fatalf("expected total issues=2 got %d", issues.Metadata.TotalData)
	}
}

// Test memastikan event FAILED benar-benar diproses oleh worker (consumer)
func TestFailedEventProcessedWithWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	x := bus.NewExchange()
	r := repo.NewInMemoryRepo()

	// Buat queue khusus test yg count event FAILED terproses
	var processed int64
	failedQ := x.Subscribe(
		constant.ExchangeTransactions,
		"test_reconcile_failed",
		func(e bus.Envelope) bool { return e.RoutingKey == constant.RKTransactionsFailed },
		64,
	)
	cons := bus.NewWorker(1, 2, failedQ, func(c context.Context, env bus.Envelope) error {
		// cast payload
		if _, ok := env.Payload.(entity.FailedTransactionOccurred); !ok {
			t.Fatalf("bad payload type: %T", env.Payload)
		}
		// simulasi kerjaan
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt64(&processed, 1)
		return nil
	})
	cons.Start(ctx)

	// start httptest server
	srv := httptest.NewServer(buildRouter(r, x))
	defer srv.Close()

	// Upload dengan 1 FAILED
	csv := `timestamp,counterparty,type,amount,status,description
1674507900, JANE, DEBIT, 125000, FAILED, atm
`
	body, ctype := multipartCSV(t, csv)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/statements", body)
	req.Header.Set("Content-Type", ctype)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("upload request error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("upload status=%d", resp.StatusCode)
	}

	// Tunggu sampai consumer memproses 1 event, atau timeout
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&processed) >= 1 {
			return // success
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("expected processed failed events >= 1, got %d", atomic.LoadInt64(&processed))
}
