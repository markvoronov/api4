package service

import (
	"context"
	"github.com/markvoronov/shortener/internal/model"
	"log/slog"
)

// ShortenerService описывает весь use-case «сокращение + разворачивание» ссылки.
type ShortenerService interface {
	// SaveOriginalUrl Save генерирует alias для оригинального URL
	// (или возвращает существующий, если уже был сохранён).
	SaveOriginalUrl(ctx context.Context, link model.ShortLink) error

	// GetOriginalUrl Get возвращает оригинальный URL по alias
	// (или ошибку, если алиас не найден).
	GetOriginalUrl(ctx context.Context, alias string) (model.ShortLink, error)

	GetAllUrls(ctx context.Context) ([]model.ShortLink, error)
}

type shortenerSvc struct {
	repo   ShortenerService
	logger *slog.Logger
}

func NewShortenerService(repo ShortenerService, logger *slog.Logger) *shortenerSvc {
	return &shortenerSvc{
		repo:   repo,
		logger: logger,
	}
}

func (s *shortenerSvc) SaveOriginalUrl(ctx context.Context, originalURL string) (string, error) {
	// проверка дубликата, генерация alias, сохранение
	s.logger.Debug("SaveOriginalUrl", slog.String("originalURL", originalURL))
	return "", nil
}

func (s *shortenerSvc) GetOriginalUrl(ctx context.Context, alias string) (string, error) {
	// просто s.repo.Get(ctx, alias)
	s.logger.Debug("GetOriginalUrl, alias :" + alias)
	return "", nil
}

func (s *shortenerSvc) GetAllUrls(ctx context.Context) ([]model.ShortLink, error) {

	op := "internal.service.shortener.GetAllUrls"
	log := s.logger.With(slog.String("op", op))

	log.Debug("GetAllUrls")

	urls, err := s.repo.GetAllUrls(ctx)
	if err != nil {
		log.Error("Can`t get all urls", slog.String("error", err.Error()))
	}
	return urls, nil
}
