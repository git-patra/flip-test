// internal/statements/usecase/handle_pending_tx_usecase.go
package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"
)

type ReviewPendingTxUsecase interface {
	Execute(ctx context.Context, evt entity.PendingTransactionOccurred) error
}
