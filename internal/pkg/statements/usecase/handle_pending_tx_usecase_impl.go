package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"

	"github.com/sirupsen/logrus"
)

type reviewPendingTxUsecase struct{}

func NewReviewPendingTxUsecase() ReviewPendingTxUsecase { return &reviewPendingTxUsecase{} }

func (u *reviewPendingTxUsecase) Execute(ctx context.Context, evt entity.PendingTransactionOccurred) error {
	logrus.Infof("TrxPendingConsumer.Execute: trx desc: %s, stats: %s, upload_id: %s", evt.Transaction.Description, evt.Transaction.Status, evt.UploadID)
	return nil
}
