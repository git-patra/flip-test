package statements

import (
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"
	"boilerplate-go/internal/pkg/statements/infrastructure/repo"
	httpiface "boilerplate-go/internal/pkg/statements/interfaces/http"
	"boilerplate-go/internal/pkg/statements/usecase"
	"context"
)

type Module struct {
	Handler *httpiface.Handler
}

func InitStatements(ctx context.Context, x *bus.Exchange) *Module {
	// infra
	memoryRepo := repo.NewInMemoryRepo()

	// usecase (producer)
	parseCSV := usecase.NewParseCSVUsecase(memoryRepo, x)
	getBalanceUsecase := usecase.NewGetBalanceUsecase(memoryRepo)
	getIssuesUsecase := usecase.NewGetIssuesUsecase(memoryRepo)

	// http handler
	h := httpiface.NewHandler(parseCSV, getBalanceUsecase, getIssuesUsecase)

	return &Module{
		Handler: h,
	}
}
