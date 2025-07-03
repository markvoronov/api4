package random

import (
	"math/rand"
	"time"
)

func newRNG() *rand.Rand {
	// Создаём новый Source, сеем его по текущему времени
	src := rand.NewSource(time.Now().UnixNano())
	// Оборачиваем в Rand — теперь можно безопасно вызывать методы
	return rand.New(src)
}

func NewRandomString(length int, r *rand.Rand) string {
	if r == nil {
		r = newRNG()
	}

	result := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, byte('a'+r.Intn(26)))
	}

	return string(result)
}
