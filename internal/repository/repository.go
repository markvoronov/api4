package repository

import (
	"context"
	"errors"
)

var (
	ErrURLNotFound    = errors.New("url not found")
	ErrURLExists      = errors.New("url exists")
	ErrAliasExists    = errors.New("alias exists")
	ErrAliasNotExists = errors.New("alias not exists")
	ErrNoConnectToDb  = errors.New("db is not connected")
)

type Repo interface {
	Get(ctx context.Context, alias string) (string, error)
	Add(ctx context.Context, url string, alias string) error
	Ping(ctx context.Context) error
}
