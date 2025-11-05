package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"

	"github.com/sirupsen/logrus"
)

type ReconcileFailedTxUsecase struct{}

func NewReconcileFailedTxUsecase() *ReconcileFailedTxUsecase { return &ReconcileFailedTxUsecase{} }

func (u *ReconcileFailedTxUsecase) Execute(ctx context.Context, evt entity.FailedTransactionOccurred) error {
	logrus.Infof("Consumer.Execute: trx desc: %s, stats: %s, upload_id: %s", evt.Transaction.Description, evt.Transaction.Status, evt.UploadID)
	return nil
}
