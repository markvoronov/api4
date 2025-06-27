package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/markvoronov/shortener/internal/config"
	"github.com/markvoronov/shortener/internal/repository/memory"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI_Save(t *testing.T) {

	tests := []struct {
		name             string
		logger           *slog.Logger
		method           string
		path             string
		body             string
		expectedCode     int
		expectedLocation string
	}{

		{
			name:             "Positive 1",
			logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
			method:           http.MethodPost,
			path:             "",
			body:             "http://ya.ru",
			expectedCode:     201,
			expectedLocation: "s",
		},
		{
			name:             "Negative 1",
			logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
			method:           http.MethodPost,
			path:             "",
			body:             "",
			expectedCode:     400,
			expectedLocation: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				BaseUrl:     "https://dzen.ru",
				AliasLength: 6,
			}
			api_test := New(cfg, tt.logger, memory.NewStorage())
			api_test.ConfigureRouterField()

			req := httptest.NewRequest(tt.method, "/"+tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "text/plain")
			writer := httptest.NewRecorder()
			rc := chi.NewRouteContext()
			rc.URLParams.Add("id", tt.path)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			api_test.RootHandle(writer, req)

			assert.Equal(t, tt.expectedCode, writer.Code, "Код ответа не совпадает с ожидаемым")
			if tt.expectedLocation != "" {
				assert.NotEmpty(t, writer.Header().Get("Location"))
			}

		})
	}
}
