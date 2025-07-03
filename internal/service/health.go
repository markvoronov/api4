package service

import "context"

// Pinger — минимальный интерфейс для нашего хранилища
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthService инкапсулирует логику «проверить доступность»
type HealthService struct {
	repo Pinger
}

func NewHealthService(repo Pinger) *HealthService {
	return &HealthService{repo: repo}
}

// Ping проверяет, доступно ли хранилище
func (s *HealthService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
