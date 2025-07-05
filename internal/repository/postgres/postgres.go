package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/markvoronov/shortener/internal/config"
	"github.com/markvoronov/shortener/internal/model"
	"github.com/markvoronov/shortener/internal/repository"
	"github.com/markvoronov/shortener/internal/service"
	"time"
)

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

	if err := s.Ping(context.Background()); err != nil {
		db.Close()      // НЕ забываем закрыть открытое соединение
		return nil, err // и возвращаем ошибку дальше
	}

	return s, nil

}

func (s *Storage) Ping(ctx context.Context) error {

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Используем PingContext вместо Ping
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

func (s *Storage) GetOriginalUrl(ctx context.Context, alias string) (model.ShortLink, error) {
	const op = "repository.postgres.GetOriginalUrl"
	var link model.ShortLink

	err := s.db.QueryRowContext(ctx, "SELECT id, url, alias, created_at FROM urls WHERE alias = $1", alias).
		Scan(&link.ID, &link.Original, &link.Alias, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return link, fmt.Errorf("%s: %w", op, repository.ErrAliasNotExists)
		} else {
			return link, fmt.Errorf("%s: %w", op, err)
		}
	}

	return link, nil
}

func (s *Storage) SaveOriginalUrl(ctx context.Context, link model.ShortLink) (model.ShortLink, error) {

	const op = "internal.repository.postgres.SaveOriginalUrl"

	// Выполняем INSERT и сразу возвращаем id и created_at
	query := `
        INSERT INTO urls (url, alias)
        VALUES ($1, $2)
        RETURNING id, created_at
        `
	row := s.db.QueryRowContext(ctx, query, link.Original, link.Alias)

	if err := row.Scan(&link.ID, &link.CreatedAt); err != nil {
		// Обрабатываем уникальные нарушения как раньше
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "urls_alias_key":
				return link, fmt.Errorf("%s: %w", op, repository.ErrAliasExists)
			case "urls_url_key":
				return link, fmt.Errorf("%s: %w", op, repository.ErrURLExists)
			}
		}
		return link, fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}

func (s *Storage) GetAllUrls(ctx context.Context) ([]model.ShortLink, error) {

	const op = "repository.postgres.GetAllUrls"

	// Подготовим запрос к нужным колонкам
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, url, alias, created_at
         FROM urls`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var urls []model.ShortLink
	for rows.Next() {
		var link model.ShortLink
		if err := rows.Scan(&link.ID, &link.Original, &link.Alias, &link.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		urls = append(urls, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return urls, nil
}

var _ service.ShortenerService = (*Storage)(nil)
