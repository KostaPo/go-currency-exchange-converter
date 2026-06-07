package exchangerate

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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	slog.InfoContext(ctx, "get all exchange rates request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
	)

	rates, err := h.svc.GetAll(ctx)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rates)
}
