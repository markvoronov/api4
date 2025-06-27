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
	"testing"
)

func TestAPI_IdPageHandle(t *testing.T) {

	tests := []struct {
		name             string
		logger           *slog.Logger
		method           string
		path             string
		expectedCode     int
		expectedLocation string
	}{

		{
			name:             "Пустой id",
			logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
			method:           http.MethodGet,
			path:             "",
			expectedCode:     400,
			expectedLocation: "",
		},
		{
			name:             "Передан корректный id",
			logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
			method:           http.MethodGet,
			path:             "some",
			expectedCode:     307,
			expectedLocation: "Location",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cfg := &config.Config{BaseUrl: "https://dzen.ru"}
			api_test := New(cfg, tt.logger, memory.NewStorage())
			api_test.ConfigureRouterField()

			req := httptest.NewRequest(tt.method, "/"+tt.path, nil)
			writer := httptest.NewRecorder()
			rc := chi.NewRouteContext()
			rc.URLParams.Add("id", tt.path)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			api_test.IdPageHandle(writer, req)

			assert.Equal(t, tt.expectedCode, writer.Code, "Код ответа не совпадает с ожидаемым")
			if tt.expectedLocation != "" {
				assert.NotEmpty(t, writer.Header().Get("Location"))
			}

		})
	}
}
