package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"

	"github.com/sirupsen/logrus"
)

type reconcileFailedTxUsecase struct{}

func NewReconcileFailedTxUsecase() ReconcileFailedTxUsecase {
	return &reconcileFailedTxUsecase{}
}

func (u *reconcileFailedTxUsecase) Execute(ctx context.Context, evt entity.FailedTransactionOccurred) error {
	logrus.Infof("TrxFailedConsumer.Execute: trx desc: %s, stats: %s, upload_id: %s", evt.Transaction.Description, evt.Transaction.Status, evt.UploadID)
	return nil
}
