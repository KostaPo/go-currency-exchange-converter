package currency

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Handler логирует HTTP-событие: какой endpoint вызван и кем.
	// request_id уже лежит в контексте — RequestID middleware положил его
	// до того как запрос дошёл сюда.
	slog.InfoContext(ctx, "get all currencies request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
		"method", r.Method,
		"path", r.URL.Path,
	)

	currencies, err := h.svc.GetAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get currencies",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "handler",
			"error", err,
		)
		http.Error(w, "failed to get currencies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(currencies)
}
