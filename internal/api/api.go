package api

import (
	"fmt"
	"github.com/markvoronov/shortener/internal/config"
	"github.com/markvoronov/shortener/session"
	"log/slog"
	"net/http"
)

type API struct {
	// Unexported field
	config          *config.Config
	logger          *slog.Logger
	router          *APIRouter
	sessionProvider *session.SessionProvider // ← добавляем
}

func New(config *config.Config, logger *slog.Logger, router *APIRouter) *API {
	return &API{
		config:          config,
		logger:          logger,
		router:          router,
		sessionProvider: &session.SessionProvider{Config: config}, // ← инициализируем
	}
}

func (a *APIRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// Start http server, configure loggers, database connection
func (api *API) Start() error {

	api.logger.Debug("Starting configure Router field")
	api.router.ConfigureRouterField()

	srv := &http.Server{
		Addr:         api.config.Address,
		Handler:      api.router,
		ReadTimeout:  api.config.HTTPServer.Timeout,
		WriteTimeout: api.config.HTTPServer.Timeout,
		IdleTimeout:  api.config.HTTPServer.IdleTimeout,
	}

	api.logger.Info("starting server", slog.String("address", api.config.Address))

	if err := srv.ListenAndServe(); err != nil {
		api.logger.Error("failed to start server")
		return fmt.Errorf("failed to start server %w", err)
	}

	return nil
}
