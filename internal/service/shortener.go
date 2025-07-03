package service

import (
	"context"
	"github.com/markvoronov/shortener/internal/repository"
)

// ShortenerService описывает весь use-case «сокращение + разворачивание» ссылки.
type ShortenerService interface {
	// Save генерирует alias для оригинального URL
	// (или возвращает существующий, если уже был сохранён).
	Save(ctx context.Context, originalURL string, alias string) (err error)

	// Get возвращает оригинальный URL по alias
	// (или ошибку, если алиас не найден).
	Get(ctx context.Context, alias string) (originalURL string, err error)
}

type shortenerSvc struct {
	repo repository.Repo
}

func NewShortenerService(repo repository.Repo) ShortenerService {
	return &shortenerSvc{repo: repo}
}

func (s *shortenerSvc) Save(ctx context.Context, originalURL string, alias string) error {
	// проверка дубликата, генерация alias, сохранение
	return nil
}

func (s *shortenerSvc) Get(ctx context.Context, alias string) (string, error) {
	// просто s.repo.Get(ctx, alias)
	return "", nil
}
