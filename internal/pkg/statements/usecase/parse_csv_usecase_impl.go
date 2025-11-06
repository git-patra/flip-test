package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/constant"
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"boilerplate-go/internal/pkg/statements/domain/response"
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"
	"boilerplate-go/internal/pkg/statements/infrastructure/repo"
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type parseCSVUsecase struct {
	repo repo.InMemoryRepo
	xchg bus.Exchange
}

func NewParseCSVUsecase(r repo.InMemoryRepo, xchg bus.Exchange) ParseCSVUsecase {
	return &parseCSVUsecase{repo: r, xchg: xchg}
}

func (u *parseCSVUsecase) Execute(ctx context.Context, r io.Reader) (result response.StatementResponse, err error) {
	result.UploadID = uuid.NewString()

	cr := csv.NewReader(bufio.NewReader(r))
	cr.FieldsPerRecord = -1
	line := 0

	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return result, err
		}

		line++
		if line == 1 { // skip header
			continue
		}
		if len(rec) < 6 {
			continue
		}

		ts, _ := strconv.ParseInt(strings.TrimSpace(rec[0]), 10, 64)
		cp := strings.TrimSpace(rec[1])
		tType := constant.TxType(strings.ToUpper(strings.TrimSpace(rec[2])))
		amt, _ := strconv.ParseInt(strings.TrimSpace(rec[3]), 10, 64)
		status := constant.TxStatus(strings.ToUpper(strings.TrimSpace(rec[4])))
		desc := strings.TrimSpace(rec[5])

		tx := entity.Transaction{
			ID: uuid.NewString(), UploadID: result.UploadID, OccurredAt: time.Unix(ts, 0),
			Counterparty: cp, Type: tType, Amount: amt, Status: status, Description: desc, Line: line,
		}
		u.repo.Save(tx)

		if tx.Status != constant.SUCCESS {
			key := result.UploadID + "/" + tx.ID
			switch status {
			case constant.FAILED:
				u.xchg.Publish(constant.ExchangeTransactions, bus.Envelope{
					RoutingKey: constant.RKTransactionsFailed,
					Key:        key,
					Payload:    entity.FailedTransactionOccurred{UploadID: result.UploadID, Transaction: tx},
				})
			case constant.PENDING:
				u.xchg.Publish(constant.ExchangeTransactions, bus.Envelope{
					RoutingKey: constant.RKTransactionsPending,
					Key:        key,
					Payload:    entity.PendingTransactionOccurred{UploadID: result.UploadID, Transaction: tx},
				})
			}
		}
	}

	return
}
