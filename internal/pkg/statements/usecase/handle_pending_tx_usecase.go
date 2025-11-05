// internal/statements/usecase/handle_pending_tx_usecase.go
package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"

	"github.com/sirupsen/logrus"
)

type ReviewPendingTxUsecase struct{}

func NewReviewPendingTxUsecase() *ReviewPendingTxUsecase { return &ReviewPendingTxUsecase{} }

func (u *ReviewPendingTxUsecase) Execute(ctx context.Context, evt entity.PendingTransactionOccurred) error {
	logrus.Infof("Consumer.Execute: trx desc: %s, stats: %s, upload_id: %s", evt.Transaction.Description, evt.Transaction.Status, evt.UploadID)
	return nil
}
