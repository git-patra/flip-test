package statements

import (
	"boilerplate-go/internal/pkg/statements/domain/constant"
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"
	"boilerplate-go/internal/pkg/statements/usecase"
	"context"
	"fmt"
)

type EventModule struct {
	Exchange bus.Exchange
}

func InitEventConsumers(ctx context.Context, x bus.Exchange) *EventModule {
	// Bind queues
	failedQ := x.Subscribe(
		constant.ExchangeTransactions,
		constant.QueueReconcileFailed,
		func(e bus.Envelope) bool { return e.RoutingKey == constant.RKTransactionsFailed },
		1024,
	)
	pendingQ := x.Subscribe(
		constant.ExchangeTransactions,
		constant.QueueReviewPending,
		func(e bus.Envelope) bool { return e.RoutingKey == constant.RKTransactionsPending },
		1024,
	)

	// Usecases
	failedUC := usecase.NewReconcileFailedTxUsecase()
	pendingUC := usecase.NewReviewPendingTxUsecase()

	// Consumers
	failedConsumer := bus.NewWorker(4, 3, failedQ, func(ctx context.Context, env bus.Envelope) error {
		p, ok := env.Payload.(entity.FailedTransactionOccurred)
		if !ok {
			return fmt.Errorf("bad payload for %s", env.RoutingKey)
		}
		return failedUC.Execute(ctx, p)
	})
	failedConsumer.Start(ctx)

	pendingConsumer := bus.NewWorker(2, 2, pendingQ, func(ctx context.Context, env bus.Envelope) error {
		p, ok := env.Payload.(entity.PendingTransactionOccurred)
		if !ok {
			return fmt.Errorf("bad payload for %s", env.RoutingKey)
		}
		return pendingUC.Execute(ctx, p)
	})
	pendingConsumer.Start(ctx)

	return &EventModule{Exchange: x}
}
