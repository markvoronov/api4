package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
)

type fullUrl struct {
	Url string `json:"url"`
}

type shortUrl struct {
	Result string `json:"result"`
}

func (api *API) ApiShortenHandle(w http.ResponseWriter, r *http.Request) {

	var fullUrl fullUrl
	var shortUrl shortUrl

	const op = "internal.api.shorten_info.ApiShortenHandle"
	log := api.logger.With(slog.String("op", op))

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		log.Debug("Method Not Allowed", slog.String("got method", r.Method))
		return
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		log.Debug("Content-Type must be application/json", slog.String("got Content-Type", mediaType))
		return
	}

	// Ограничиваем размер тела
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Body read error: %v", err), http.StatusBadRequest)
		log.Debug("Body read error", slog.Any("Error", err))
		return
	}
	defer r.Body.Close()

	// десериализуем JSON в fullUrl
	if err = json.Unmarshal(body, &fullUrl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Debug("Body read error", slog.Any("Error", err))
		return
	}

	if fullUrl.Url == "" {
		http.Error(w, "Field 'url' is required", http.StatusBadRequest)
		log.Debug("Field 'url' is required")
		return
	}
	//fmt.Println(fullUrl.Url)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Пишем тело и проверяем ошибку

	shortUrl.Result = api.config.BaseUrl

	//Кодируем в JSON
	out, err := json.Marshal(shortUrl)
	if err != nil {
		log.Info("Unable marshall json", slog.Any("error", err))
		return
	}

	if _, err := w.Write(out); err != nil {
		//log.Printf("Failed to write response: %v", err)
		log.Info("Unable to write answer", slog.Any("error", err))
		return
	}

}
