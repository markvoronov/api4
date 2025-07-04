package model

import "time"

type ShortLink struct {
	ID        int64     // первичный ключ в БД
	Alias     string    // собственно, сокращённая часть
	Original  string    // оригинальный URL
	CreatedAt time.Time // время создания
	CreatedBy string    // ID пользователя, если нужно
	// ...потом сюда можно добавить любые новые поля
}

// Используется при запросах по пути /shorten
type UrlProcessingQ struct {
	URL string `json:"url"`
}

// Используется при ответах по пути /shorten
type UrlProcessingA struct {
	URL string `json:"result"`
}
