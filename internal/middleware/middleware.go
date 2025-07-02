package middleware

import (
	"context"
	"github.com/markvoronov/shortener/session"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(sp *session.SessionProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := sp.ParseToken(r)
			if err != nil {
				// 1) Если токена нет или он неверный – создаём новый userID
				//TODO : получать из БД
				userID = 1 /* здесь ваша логика: uuid.New().String() или из БД */

				// 2) и тут же устанавливаем куку с этим новым userID
				if err := sp.GenerateTokenAndSetCookie(w, userID); err != nil {
					http.Error(w, "failed to set cookie", http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
