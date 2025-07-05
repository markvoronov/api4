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
