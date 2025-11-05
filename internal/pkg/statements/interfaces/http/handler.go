package httpiface

import (
	"boilerplate-go/internal/pkg/statements/usecase"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	parseCSV   *usecase.ParseCSVUsecase
	getBalance *usecase.GetBalanceUsecase
	getIssues  *usecase.GetIssuesUsecase
}

func NewHandler(p *usecase.ParseCSVUsecase, g *usecase.GetBalanceUsecase, i *usecase.GetIssuesUsecase) *Handler {
	return &Handler{parseCSV: p, getBalance: g, getIssues: i}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	id, balance, issues, err := h.parseCSV.Execute(r.Context(), file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"upload_id": id, "balance": balance, "issues": issues,
	})
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("upload_id")
	bal := h.getBalance.Execute(id)
	json.NewEncoder(w).Encode(map[string]any{"upload_id": id, "balance": bal})
}

func (h *Handler) GetIssues(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("upload_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	statuses := []string{}
	if s := r.URL.Query().Get("status"); s != "" {
		statuses = strings.Split(s, ",")
	}

	items, total := h.getIssues.Execute(id, statuses, page, size)
	json.NewEncoder(w).Encode(map[string]any{
		"upload_id": id,
		"total":     total,
		"page":      page,
		"size":      size,
		"data":      items,
	})
}
