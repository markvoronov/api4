package memory

import (
	"fmt"
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

func (s *Storage) Get(alias string) (string, error) {

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

func (s *Storage) Add(url string, alias string) error {

	if url == "" {
		return fmt.Errorf("url empty")
	}
	if alias == "" {
		return fmt.Errorf("alias empty")
	}

	// Проверим, есть ли такой алиас
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if m, ok := s.Store[alias]; ok {

		return fmt.Errorf("alias already exists for %s: %w", m, repository.ErrAliasExists)
	}

	// Проверим, есть ли такой url
	var findUrl bool
	for _, val := range s.Store {
		if val == url {
			findUrl = true
		}
	}
	if findUrl {

		return fmt.Errorf("url already exists %s: %w", url, repository.ErrURLExists)
	}

	s.Store[alias] = url

	return nil

}
