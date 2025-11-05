package statements

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"
	"boilerplate-go/internal/pkg/statements/infrastructure/repo"
	httpiface "boilerplate-go/internal/pkg/statements/interfaces/http"
	"boilerplate-go/internal/pkg/statements/usecase"
	"context"
)

type Module struct {
	Handler *httpiface.Handler
	Bus     *bus.InMemoryBus
}

func InitStatements(ctx context.Context) *Module {
	repository := repo.NewInMemoryRepo()
	eventBus := bus.NewInMemoryBus(1024)

	// start worker pool
	reconHandler := func(ctx context.Context, evt entity.FailedTransactionOccurred) error {
		// bisa logika reconciliation atau dummy
		return nil
	}
	consumer := bus.NewWorkerPool(4, 3, eventBus.Subscribe(), reconHandler)
	consumer.Start(ctx)

	parseCSV := usecase.NewParseCSVUsecase(repository, eventBus)
	getBalance := usecase.NewGetBalanceUsecase(repository)
	getIssues := usecase.NewGetIssuesUsecase(repository)
	handler := httpiface.NewHandler(parseCSV, getBalance, getIssues)

	return &Module{Handler: handler, Bus: eventBus}
}
