package currency

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"encoding/json"
	"errors"
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

	slog.InfoContext(ctx, "get all currencies request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
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

func (h *Handler) GetByCode(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	slog.InfoContext(ctx, "get currency by code request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"code", code,
	)

	currency, err := h.svc.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "currency not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(ctx, "failed to get currency",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "handler",
			"code", code,
			"error", err,
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(currency)
}
