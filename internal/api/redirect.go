package api

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/markvoronov/shortener/internal/repository"

	//"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func (api *API) IdPageHandle(w http.ResponseWriter, r *http.Request) {
	const op = "internal.api.redirect.IdPageHandle"
	log := api.logger.With(slog.String("op", op))

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Debug("Method Not Allowed", slog.String("Get method", r.Method))
		return
	}

	idPage := chi.URLParam(r, "id")
	if idPage == "" {
		log.Info("id Page is empty")
		w.WriteHeader(http.StatusBadRequest)
		//render.JSON(writer, request, resp.Error("invalid request"))

		return
	}

	log.Debug("alias was getting", slog.String("idPage", idPage))

	baseUrl, err := api.storage.Get(idPage)
	if errors.Is(err, repository.ErrAliasNotExists) {
		log.Debug("not found alias in DB", slog.String("alias", idPage))
		w.Write([]byte("not found alias " + idPage + " in DB"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Debug("error while getting url by alias", slog.Any("error", err.Error()))
		w.Write([]byte("something error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", baseUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
	log.Debug("redirect to", slog.String("originURL", baseUrl))
	// redirect to found url
	http.Redirect(w, r, baseUrl, http.StatusFound)

}

type URLGetter interface {
	GetURL(alias string) (string, error)
}
