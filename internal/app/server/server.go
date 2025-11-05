package server

import (
	customMiddleware "boilerplate-go/api/middleware"
	"boilerplate-go/api/routes"
	"boilerplate-go/config"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	chiRouter *chi.Mux
	config    *config.AppConfig
}

func NewServer(
	cfg *config.AppConfig,
) *Server {

	return &Server{
		chiRouter: initializeChiRouter(
			cfg,
		),
		config: cfg,
	}
}

func (server *Server) Start() error {
	httpServer := http.Server{
		Addr:        fmt.Sprintf(":%d", server.config.ServerPort),
		Handler:     server.chiRouter,
		ReadTimeout: server.config.RequestTimeout,
	}

	err := ServeHTTP(&httpServer, httpServer.Addr, 0)
	if err != nil {
		logrus.Error("failed to start the REST API server:", err)
		return err
	}

	logrus.Info("REST API server stopped")
	return nil
}

func initializeChiRouter(
	config *config.AppConfig,
) *chi.Mux {
	chiRouter := chi.NewRouter()

	// Middlewares.
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(customMiddleware.CORSMiddleware)

	// Routes.
	routes.RegisterRoutes(chiRouter, config)

	return chiRouter
}
