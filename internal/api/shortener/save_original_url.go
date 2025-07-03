package shortener

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"time"
)

//
//// Интерфейс объявляем здесь, рядом с handler’ом, потому что именно он его вызывает.
//type SaveService interface {
//	Save(ctx context.Context, url string, alias string) error
//}

func (h *Handler) SaveOriginalUrl(w http.ResponseWriter, r *http.Request) {
	const op = "internal.api.save.RootHandle"
	log := h.logger.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		h.logger.Debug("Method Not Allowed", slog.String("Get method", r.Method))
		return
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusUnsupportedMediaType)
		log.Debug("Content-Type must be text/plain", slog.String("Got Content-Type", contentType))
		return
	}

	// Ограничиваем размер тела
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Body read error: %v", err), http.StatusBadRequest)
		log.Debug("Body read error", slog.Any("error", err))
		return
	}
	defer r.Body.Close()

	bodyStr := string(body)
	if bodyStr == "" {
		http.Error(w, "Body empty", http.StatusBadRequest)
		log.Debug("Body empty")
		return
	}

	alias, err := h.service.SaveOriginalUrl(ctx, bodyStr)

	if err != nil {
		log.Info("Error while save url", slog.String("url", err.Error()))
		w.Write([]byte("Error while save url"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", bodyStr)
	w.WriteHeader(http.StatusCreated)
	// Пишем тело и проверяем ошибку
	//	ref := api.config.Address + "/" + alias
	if _, err := w.Write([]byte(alias)); err != nil {
		//log.Printf("Failed to write response: %v", err)
		return
	}
	log.Info("Записан адрес", slog.String("alias", alias))

}
