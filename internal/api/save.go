package api

import (
	"errors"
	"fmt"
	"github.com/markvoronov/shortener/internal/api/random"
	"github.com/markvoronov/shortener/internal/repository"
	"io"
	"log/slog"
	"mime"
	"net/http"
)

func (api *API) RootHandle(w http.ResponseWriter, r *http.Request) {
	const op = "internal.api.save.RootHandle"
	log := api.logger.With(slog.String("op", op))

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		log.Debug("Method Not Allowed", slog.String("Get method", r.Method))
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

	alias := random.NewRandomString(api.config.AliasLength, nil)
	err = api.storage.Add(bodyStr, alias)
	if errors.Is(err, repository.ErrURLExists) {
		log.Info("url already exists", slog.String("url", bodyStr))
		w.Write([]byte("url already exists"))
		//render.JSON(w, r, resp.Error("url already exists"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", bodyStr)
	w.WriteHeader(http.StatusCreated)
	// Пишем тело и проверяем ошибку
	ref := api.config.Address + "/" + alias
	if _, err := w.Write([]byte(ref)); err != nil {
		//log.Printf("Failed to write response: %v", err)
		return
	}
	log.Info("Записан адрес", slog.String("alias", alias), slog.String("ref", ref))

}
