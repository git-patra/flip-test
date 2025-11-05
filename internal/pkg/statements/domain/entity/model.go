package entity

import (
	"boilerplate-go/internal/pkg/statements/domain/constant"
	"time"
)

type Transaction struct {
	ID           string            `json:"id"`
	UploadID     string            `json:"upload_id"`
	OccurredAt   time.Time         `json:"occurred_at"`
	Counterparty string            `json:"counterparty"`
	Type         constant.TxType   `json:"type"`
	Amount       int64             `json:"amount"`
	Status       constant.TxStatus `json:"status"`
	Description  string            `json:"description"`
	Line         int               `json:"line"`
}

type FailedTransactionOccurred struct {
	UploadID    string      `json:"upload_id"`
	Transaction Transaction `json:"transaction"`
}

type PendingTransactionOccurred struct {
	UploadID    string
	Transaction Transaction `json:"transaction"`
}
