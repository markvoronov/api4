package repository

import (
	"errors"
)

var (
	ErrURLNotFound    = errors.New("url not found")
	ErrURLExists      = errors.New("url exists")
	ErrAliasExists    = errors.New("alias exists")
	ErrAliasNotExists = errors.New("alias not exists")
	ErrNoConnectToDb  = errors.New("db is not connected")
)

//
//type Storage interface {
//	GetOriginalUrl(ctx context.Context, alias string) (string, error)
//	SaveOriginalUrl(ctx context.Context, model model.ShortLink) error
//	Ping(ctx context.Context) error
//	GetAllUrls(ctx context.Context) ([]model.ShortLink, error)
//}
