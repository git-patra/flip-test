package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/constant"
	"boilerplate-go/internal/pkg/statements/domain/entity"
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

type ParseCSVUsecase struct {
	repo *repo.InMemoryRepo
	bus  *bus.InMemoryBus
}

func NewParseCSVUsecase(r *repo.InMemoryRepo, b *bus.InMemoryBus) *ParseCSVUsecase {
	return &ParseCSVUsecase{repo: r, bus: b}
}

func (u *ParseCSVUsecase) Execute(ctx context.Context, r io.Reader) (uploadID string, balance int64, issues int, err error) {
	uploadID = uuid.NewString()

	cr := csv.NewReader(bufio.NewReader(r))
	cr.FieldsPerRecord = -1
	line := 0

	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", 0, 0, err
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
			ID: uuid.NewString(), UploadID: uploadID, OccurredAt: time.Unix(ts, 0),
			Counterparty: cp, Type: tType, Amount: amt, Status: status, Description: desc, Line: line,
		}
		u.repo.Save(tx)

		if tx.Status == constant.SUCCESS {
			if tx.Type == constant.CREDIT {
				balance += tx.Amount
			} else if tx.Type == constant.DEBIT {
				balance -= tx.Amount
			}
		} else if tx.Status == constant.FAILED || tx.Status == constant.PENDING {
			issues++
		}

		if tx.Status == constant.FAILED {
			u.bus.PublishNonBlocking(entity.FailedTransactionOccurred{UploadID: uploadID, TxID: tx.ID})
		}
	}

	return uploadID, balance, issues, nil
}
