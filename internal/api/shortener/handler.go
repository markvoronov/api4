package shortener

import (
	"context"
	"github.com/markvoronov/shortener/internal/model"
	"log/slog"
)

type Service interface {
	GetOriginalUrl(ctx context.Context, alias string) (string, error)
	SaveOriginalUrl(ctx context.Context, url string) (string, error)
	GetAllUrls(ctx context.Context) ([]model.ShortLink, error)
}

type Handler struct {
	service Service // интерфейс, который умеет сохранять url и получать их по алиасу
	logger  *slog.Logger
}

func NewHandler(service Service, log *slog.Logger) *Handler {
	return &Handler{service: service, logger: log}
}
