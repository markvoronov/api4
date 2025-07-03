package shortener

import (
	"encoding/json"
	"fmt"
	"github.com/markvoronov/shortener/internal/api"
	authmw "github.com/markvoronov/shortener/internal/middleware"
	"github.com/markvoronov/shortener/internal/repository"
	"github.com/markvoronov/shortener/pkg/random"
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

func (api *api.API) ApiShortenHandle(w http.ResponseWriter, r *http.Request) {

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

	_, ok := r.Context().Value(authmw.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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

	alias := random.NewRandomString(api.config.AliasLength, nil)
	err = api.storage.Add(bodyStr, alias)
	if errors.Is(err, repository.ErrURLExists) {
		log.Info("url already exists", slog.String("url", bodyStr))
		w.Write([]byte("url already exists"))
		//render.JSON(w, r, resp.Error("url already exists"))
		return
	}

	if err != nil {
		log.Info("Error while add url", slog.String("url", err.Error()))
		w.Write([]byte("Error while add url"))
		//render.JSON(w, r, resp.Error("url already exists"))
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
