package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markvoronov/shortener/internal/config"
	mwLogger "github.com/markvoronov/shortener/internal/middleware/logger"
	"github.com/markvoronov/shortener/internal/repository"

	//"github.com/markvoronov/shortener/internal/handler/redirect"
	//	"github.com/markvoronov/shortener/internal/handler/save"
	"log/slog"
	"net/http"
)

type API struct {
	// Unexported field
	config  *config.Config
	logger  *slog.Logger
	router  *chi.Mux
	storage repository.Storage
}

func New(config *config.Config, logger *slog.Logger, storage repository.Storage) *API {
	return &API{
		config:  config,
		logger:  logger,
		router:  chi.NewRouter(),
		storage: storage,
	}
}

func (api *API) ConfigureRouterField() {

	api.router.Use(middleware.Recoverer)
	api.router.Use(middleware.RequestID)
	api.router.Use(middleware.RealIP)
	api.router.Use(middleware.URLFormat)
	api.router.Use(mwLogger.New(api.logger))

	api.router.Post("/", api.RootHandle)
	api.router.Get("/", api.IdPageHandle)
	api.router.Get("/{id}", api.IdPageHandle)
	api.router.Post("/shorten", api.ApiShortenHandle)
	api.router.Get("/ping", api.PingHandle)

	// Обработчик для несуществующих маршрутов
	api.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 — страница не найдена (Mark)"))
	})

}

// Start http server, configure loggers, database connection
func (api *API) Start() error {

	api.logger.Debug("Starting configure Router field")
	api.ConfigureRouterField()

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
