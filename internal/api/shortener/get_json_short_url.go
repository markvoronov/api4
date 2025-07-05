package shortener

import (
	"context"
	"encoding/json"
	"github.com/markvoronov/shortener/internal/model"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) GetJSONShortUrl(writer http.ResponseWriter, request *http.Request) {
	const op = "internal.api.shortener.get_json_short_url.GetJSONShortUrl"
	log := h.logger.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancel()

	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", "POST")
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Debug("Method Not Allowed", slog.String("Get method", request.Method))
		return
	}

	var urlProcessing model.UrlProcessingQ

	err := json.NewDecoder(request.Body).Decode(&urlProcessing)
	if err != nil {
		log.Info("Error while decode json", slog.String("error", err.Error()))
		http.Error(writer, "Error while decode json", http.StatusBadRequest)
		return
	}

	originalUrl := urlProcessing.URL
	log.Debug("got url", slog.String("url", originalUrl))
	link, err := h.service.SaveOriginalUrl(ctx, originalUrl)

	if err != nil {
		log.Info("Error while save url", slog.String("url", err.Error()))
		http.Error(writer, "Error while save url", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	//writer.Header().Set("Location", originalUrl)
	writer.WriteHeader(http.StatusCreated)
	// Пишем тело и проверяем ошибку
	//	ref := api.config.Address + "/" + alias

	rezult := model.UrlProcessingA{
		URL: link.Alias,
	}

	err = json.NewEncoder(writer).Encode(rezult)
	if err != nil {
		log.Info("Error while encode json", slog.String("url", err.Error()))
		http.Error(writer, "Error while save url", http.StatusInternalServerError)
		return
	}

	log.Debug("Записан адрес", slog.String("alias", link.Alias))

}
