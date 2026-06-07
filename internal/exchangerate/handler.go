package exchangerate

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
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

func (h *Handler) GetByPair(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	pair := r.PathValue("pair")
	if len(pair) != 6 {
		http.Error(w, "invalid currency pair, expected format: USDRUB", http.StatusBadRequest)
		return
	}

	baseCode := strings.ToUpper(pair[:3])
	targetCode := strings.ToUpper(pair[3:])

	slog.InfoContext(ctx, "get exchange rate by pair request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
		"base", baseCode,
		"target", targetCode,
	)

	rate, err := h.svc.GetByPair(ctx, baseCode, targetCode)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "exchange rate not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(ctx, "failed to get exchange rate",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "handler",
			"base", baseCode,
			"target", targetCode,
			"error", err,
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rate)
}
