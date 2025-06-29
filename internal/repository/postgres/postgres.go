package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/markvoronov/shortener/internal/config"
	"github.com/markvoronov/shortener/internal/repository"
	"time"
)

//type Config struct {
//	Host     string
//	Port     string
//	Username string
//	Password string
//	DBName   string
//	SSLMode  string
//}

type Storage struct {
	db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*Storage, error) {

	//ps := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	//	cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := sql.Open("pgx", cfg.Database.DSN)

	if err != nil {
		return nil, err
	}

	// Устанавливаем таймауты
	db.SetConnMaxLifetime(cfg.Database.Pool.ConnLifetime)
	db.SetMaxOpenConns(cfg.Database.Pool.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.Pool.MaxIdleConns)

	s := &Storage{db: db}

	err = s.TestPing()
	if err != nil {
		return nil, err
	}

	return s, nil

}

func (s *Storage) TestPing() error {

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Используем PingContext вместо Ping
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

func (s *Storage) Get(alias string) (string, error) {
	const op = "repository.postgres.Get"
	var destUrl string

	err := s.db.QueryRow("SELECT url FROM url5 WHERE alias = $1", alias).Scan(&destUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, repository.ErrAliasNotExists)
		} else {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return destUrl, nil
}

func (s *Storage) Add(urlToSave string, alias string) error {

	const op = "internal.repository.postgres.Add"

	// В PostgreSQL можно сразу вызвать Exec без Prepare —
	// драйвер сам кэширует запросы.
	_, err := s.db.Exec(
		`INSERT INTO url5 (url, alias) VALUES ($1, $2)`,
		urlToSave,
		alias,
	)
	if err != nil {
		// 23505 — стандартный SQL-код PostgreSQL для unique_violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, repository.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
