package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"context"
)

type ReconcileFailedTxUsecase interface {
	Execute(ctx context.Context, evt entity.FailedTransactionOccurred) error
}
