package routes

import (
	"boilerplate-go/config"
	"boilerplate-go/internal/delivery/rest/response"
	"boilerplate-go/internal/pkg/statements"
	"boilerplate-go/internal/pkg/statements/infrastructure/bus"
	httpiface "boilerplate-go/internal/pkg/statements/interfaces/http"
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes configures the main routes for the application.
func RegisterRoutes(apiV1Router *chi.Mux, cfg *config.AppConfig, x *bus.Exchange) {
	// Create a sub-router for API version 1
	apiV1Router.Mount("/api/v1", apiV1Router)

	// Create a sub-router for Loans
	apiV1WithAuthRouter := chi.NewRouter()

	apiV1Router.Mount("/", apiV1WithAuthRouter)
	apiV1Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := response.BuildSuccessResponseWithData(response.Ok, "ok")
		response.JSON(w, httpResponse.StatusCode, httpResponse)
	})

	ctx := context.Background()
	mod := statements.InitStatements(ctx, x)

	apiV1WithAuthRouter.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	httpiface.RegisterRoutes(apiV1WithAuthRouter, mod.Handler)
}
