package redirect

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func IdPageHandle(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.redirect.IdPageHandle"
	log = log.With(slog.String("op", op))

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idPage := chi.URLParam(r, "id")
	if idPage == "" {
		log.Info("id Page is empty")

		//render.JSON(writer, request, resp.Error("invalid request"))

		return
	}

	log.Info("alias was getting", idPage)

	originUrl := "http://q2n2olm25c.com/hgfkcdbm6j" //Временно
	w.Header().Set("Location", originUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
	// redirect to found url
	http.Redirect(w, r, originUrl, http.StatusFound)

}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, URLGetter URLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)

		if request.Method != http.MethodGet {
			writer.Header().Set("Allow", "GET")
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		idPage := chi.URLParam(request, "id")
		if idPage == "" {
			log.Info("id Page is empty")

			//render.JSON(writer, request, resp.Error("invalid request"))

			return
		}

		log.Info("alias was getting", slog.String("idPage", idPage))

		originUrl := "http://q2n2olm25c.com/hgfkcdbm6j" //Временно
		writer.Header().Set("Location", originUrl)
		writer.WriteHeader(http.StatusTemporaryRedirect)
		// redirect to found url
		http.Redirect(writer, request, originUrl, http.StatusFound)
	}
}
