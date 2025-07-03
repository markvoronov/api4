package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markvoronov/shortener/internal/api/health"
	"github.com/markvoronov/shortener/internal/api/shortener"
	mwLogger "github.com/markvoronov/shortener/internal/middleware/logger"
	"log/slog"
	"net/http"
	"time"
)

type APIRouter struct {
	router   *chi.Mux
	logger   *slog.Logger
	healthH  *health.Handler
	shortenH *shortener.Handler
}

func NewAPIRouter(r *chi.Mux, l *slog.Logger, healthH *health.Handler, shortenH *shortener.Handler) *APIRouter {
	return &APIRouter{
		router:   r,
		logger:   l,
		healthH:  healthH,
		shortenH: shortenH,
	}
}

func (a APIRouter) ConfigureRouterField() {

	router := a.router

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(a.logger))
	router.Use(middleware.Timeout(2 * time.Second))

	router.Get("/ping", a.healthH.Ping)

	router.Post("/", a.shortenH.SaveOriginalUrl)
	//router.Get("/", a.shortenH.IdPageHandle)
	router.Get("/{id}", a.shortenH.Redirect)
	//	router.With(authmw.AuthMiddleware(api.sessionProvider)).Post("/shorten", shortener.ApiShortenHandle)
	router.Get("/ping", a.healthH.Ping)

	// Обработчик для несуществующих маршрутов
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 — страница не найдена (Mark)"))
	})

}
