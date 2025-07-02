package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/markvoronov/shortener/internal/config"
)

type SessionProvider struct {
	*config.Config
}

const (
	cookieName = "user_token"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

// Генерация JWT и установка в cookie
func (p *SessionProvider) GenerateTokenAndSetCookie(w http.ResponseWriter, userID int64) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		UserID: userID,
	})

	signed, err := token.SignedString([]byte(p.SecretKey))
	if err != nil {
		return fmt.Errorf("failed to sign token: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // или true в проде
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(1 * time.Hour),
	})
	return nil
}

// Получение JWT из cookie и его парсинг
func (p *SessionProvider) ParseToken(r *http.Request) (int64, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		return 0, fmt.Errorf("failed to get cookie: %w", err)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(p.SecretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, jwt.ErrTokenInvalidIssuer
	}
	if claims.UserID == 0 {
		return 0, jwt.ErrInvalidKey
	}

	return claims.UserID, nil
}
