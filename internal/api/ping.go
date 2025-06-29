package api

import (
	//"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func (api *API) PingHandle(w http.ResponseWriter, r *http.Request) {
	const op = "internal.api.ping.PingHandle"
	log := api.logger.With(slog.String("op", op))

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Debug("Method Not Allowed", slog.String("Get method", r.Method))
		return
	}

	err := api.storage.TestPing()
	if err != nil {
		log.Debug("error ping: %w", err)
		w.Write([]byte("error ping"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

//
//type URLGetter interface {
//	GetURL(alias string) (string, error)
//}
