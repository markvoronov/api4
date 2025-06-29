package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/markvoronov/shortener/internal/config"
	"log/slog"
)

func RunMigrations(cfg *config.Config, logger *slog.Logger) error {
	m, err := migrate.New(
		"file://migrations",
		cfg.Database.DSN,
	)
	if err != nil {
		return fmt.Errorf("create migrate: %w", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}
	logger.Info("migrations applied")
	return nil
}
