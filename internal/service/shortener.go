package service

import (
	"context"
	"github.com/markvoronov/shortener/internal/repository"
)

// ShortenerService описывает весь use-case «сокращение + разворачивание» ссылки.
type ShortenerService interface {
	// SaveOriginalUrl Save генерирует alias для оригинального URL
	// (или возвращает существующий, если уже был сохранён).
	SaveOriginalUrl(ctx context.Context, originalURL string, alias string) (err error)

	// GetOriginalUrl Get возвращает оригинальный URL по alias
	// (или ошибку, если алиас не найден).
	GetOriginalUrl(ctx context.Context, alias string) (originalURL string, err error)
}

type shortenerSvc struct {
	repo repository.Repo
}

func NewShortenerService(repo repository.Repo) ShortenerService {
	return &shortenerSvc{repo: repo}
}

func (s *shortenerSvc) SaveOriginalUrl(ctx context.Context, originalURL string, alias string) error {
	// проверка дубликата, генерация alias, сохранение
	return nil
}

func (s *shortenerSvc) GetOriginalUrl(ctx context.Context, alias string) (string, error) {
	// просто s.repo.Get(ctx, alias)
	return "", nil
}
