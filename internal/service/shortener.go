package service

import (
	"context"
	"fmt"
	"github.com/markvoronov/shortener/internal/api/shortener"
	"github.com/markvoronov/shortener/internal/config"
	"github.com/markvoronov/shortener/internal/model"
	"github.com/markvoronov/shortener/pkg/random"
	"log/slog"
)

// ShortenerService описывает весь use-case «сокращение + разворачивание» ссылки.
type ShortenerService interface {
	// SaveOriginalUrl Save генерирует alias для оригинального URL
	// (или возвращает существующий, если уже был сохранён).
	SaveOriginalUrl(ctx context.Context, link model.ShortLink) (model.ShortLink, error)

	// GetOriginalUrl Get возвращает оригинальный URL по alias
	// (или ошибку, если алиас не найден).
	GetOriginalUrl(ctx context.Context, alias string) (model.ShortLink, error)

	GetAllUrls(ctx context.Context) ([]model.ShortLink, error)

	Ping(ctx context.Context) error
}

type ShortenerSvc struct {
	repo   ShortenerService
	logger *slog.Logger
	config *config.Config
}

func NewShortenerService(repo ShortenerService, logger *slog.Logger, config *config.Config) *ShortenerSvc {
	return &ShortenerSvc{
		repo:   repo,
		logger: logger,
		config: config,
	}
}

func (s *ShortenerSvc) SaveOriginalUrl(ctx context.Context, originalURL string) (model.ShortLink, error) {
	// проверка дубликата, генерация alias, сохранение
	const op = "internal.api.save.RootHandle"
	log := s.logger.With(slog.String("op", op))
	log.Debug("SaveOriginalUrl", slog.String("originalURL", originalURL))
	alias := random.NewRandomString(s.config.AliasLength, nil)
	link := model.ShortLink{
		Original: originalURL,
		Alias:    alias,
	}

	link, err := s.repo.SaveOriginalUrl(ctx, link)
	if err != nil {
		log.Error("Can`t save new url", slog.String("error", err.Error()))
		return link, err
	}

	return link, nil
}

func (s *ShortenerSvc) GetOriginalUrl(ctx context.Context, alias string) (string, error) {
	// просто s.repo.Get(ctx, alias)
	op := "internal.service.shortener.GetAllUrls"
	log := s.logger.With(slog.String("op", op))
	log.Debug("GetOriginalUrl, alias :" + alias)

	link, err := s.repo.GetOriginalUrl(ctx, alias)
	if err != nil {
		log.Error("Can`t get original url", slog.String("error", err.Error()))
		return "", fmt.Errorf("Can`t get original url %w", err)
	}
	return "https://" + link.Original, nil
}

func (s *ShortenerSvc) GetAllUrls(ctx context.Context) ([]model.ShortLink, error) {

	op := "internal.service.shortener.GetAllUrls"
	log := s.logger.With(slog.String("op", op))

	log.Debug("GetAllUrls")

	urls, err := s.repo.GetAllUrls(ctx)
	if err != nil {
		log.Error("Can`t get all urls", slog.String("error", err.Error()))
	}
	return urls, nil
}

func NewHealthService(repo ShortenerService) *ShortenerSvc {
	return &ShortenerSvc{repo: repo}
}

// Ping проверяет, доступно ли хранилище
func (s *ShortenerSvc) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

var _ shortener.Service = (*ShortenerSvc)(nil)
