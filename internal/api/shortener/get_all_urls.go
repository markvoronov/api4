package shortener

import (
	"context"
	"encoding/json"
	"github.com/markvoronov/shortener/internal/api/response"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) GetAllUrls(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	urls, err := h.service.GetAllUrls(ctx)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		h.logger.Error("GetAllUrls failed", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response.Error("can’t get all urls"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		h.logger.Error("ailed to encode response", slog.String("err", err.Error()))
		// на случай, если сериализация всё же провалится
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
