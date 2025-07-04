package memory

import (
	"context"
	"fmt"
	"github.com/markvoronov/shortener/internal/model"
	"github.com/markvoronov/shortener/internal/repository"
	"sync"
)

type Storage struct {
	Store map[string]string
	mtx   *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		Store: make(map[string]string),
		mtx:   &sync.Mutex{},
	}
}

func (s *Storage) GetOriginalUrl(ctx context.Context, alias string) (string, error) {

	if alias == "" {
		return "", fmt.Errorf("alias empty")
	}

	// Проверим, есть ли такой алиас
	s.mtx.Lock()
	defer s.mtx.Unlock()
	m, ok := s.Store[alias]

	if !ok {
		return "", fmt.Errorf("Can`t find alias %s in DB: %w", m, repository.ErrAliasNotExists)
	}

	return m, nil

}

func (s *Storage) SaveOriginalUrl(ctx context.Context, link model.ShortLink) error {

	if link.Original == "" {
		return fmt.Errorf("url empty")
	}
	if link.Alias == "" {
		return fmt.Errorf("alias empty")
	}

	// Проверим, есть ли такой алиас
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if m, ok := s.Store[link.Alias]; ok {

		return fmt.Errorf("alias already exists for %s: %w", m, repository.ErrAliasExists)
	}

	// Проверим, есть ли такой url
	var findUrl bool
	for _, val := range s.Store {
		if val == link.Original {
			findUrl = true
		}
	}
	if findUrl {

		return fmt.Errorf("url already exists %s: %w", link.Original, repository.ErrURLExists)
	}

	s.Store[link.Alias] = link.Original

	return nil

}

func (s *Storage) Ping(ctx context.Context) error {
	return nil
}

func (s *Storage) GetAllUrls(ctx context.Context) ([]model.ShortLink, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	var urls []model.ShortLink
	for i, v := range s.Store {
		var link model.ShortLink
		link.Original = v
		link.Alias = i
		urls = append(urls, link)
	}

	return urls, nil
}
