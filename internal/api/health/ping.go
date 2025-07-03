package health

import (
	"context"
	"log/slog"
	"net/http"
)

// PingService описывает бизнес-логику проверки «здоровья» хранилища.
type PingService interface {
	// Ping выполняет проверку доступности и возвращает ошибку, если она недоступна.
	Ping(ctx context.Context) error
}

// Handler обрабатывает HTTP-эндпоинты «здоровья» (health-check).
type Handler struct {
	pingSvc PingService
	logger  *slog.Logger
}

// NewHandler создаёт новый экземпляр Handler.
// pingSvc  — реализация PingService,
// logger   — slog-логгер для записи событий.
func NewHandler(pingSvc PingService, log *slog.Logger) *Handler {
	return &Handler{pingSvc: pingSvc, logger: log}
}

// Ping обрабатывает GET /ping.
//   - Если хранилище доступно, возвращает 204 No Content.
//   - Если ping() вернёт ошибку — 500 Internal Server Error и тело "error ping".
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {

	const op = "internal.api.health.ping.Ping"
	log := h.logger.With(slog.String("op", op))

	// В chi мы уже смонтировали только GET, так что метод тут проверять не обязательно.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Debug("method not allowed", slog.String("Get method", r.Method))
		return
	}

	if err := h.pingSvc.Ping(r.Context()); err != nil {
		log.Error("ping failed", slog.String("error", err.Error()))
		http.Error(w, "error ping", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
