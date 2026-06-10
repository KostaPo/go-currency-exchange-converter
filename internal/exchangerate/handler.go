package exchangerate

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
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
		writeError(w, "internal server error", http.StatusInternalServerError)
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
		writeError(w, "invalid currency pair, expected format: USDRUB", http.StatusBadRequest)
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
			writeError(w, "exchange rate not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(ctx, "failed to get exchange rate",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "handler",
			"base", baseCode,
			"target", targetCode,
			"error", err,
		)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rate)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	baseCode := strings.ToUpper(r.FormValue("baseCurrencyCode"))
	targetCode := strings.ToUpper(r.FormValue("targetCurrencyCode"))
	rateStr := r.FormValue("rate")

	if baseCode == "" || targetCode == "" || rateStr == "" {
		writeError(w, "baseCurrencyCode, targetCurrencyCode and rate are required", http.StatusBadRequest)
		return
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		writeError(w, "invalid rate value", http.StatusBadRequest)
		return
	}

	slog.InfoContext(ctx, "create exchange rate request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
		"base", baseCode,
		"target", targetCode,
		"rate", rate,
	)

	er, err := h.svc.Create(ctx, baseCode, targetCode, rate)
	if err != nil {
		switch {
		case errors.Is(err, ErrAlreadyExists):
			writeError(w, "exchange rate already exists", http.StatusConflict)
		case strings.Contains(err.Error(), "not found"):
			writeError(w, err.Error(), http.StatusNotFound)
		default:
			slog.ErrorContext(ctx, "failed to create exchange rate",
				"request_id", middleware.IDFromContext(ctx),
				"layer", "handler",
				"error", err,
			)
			writeError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(er)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	pair := r.PathValue("pair")
	if len(pair) != 6 {
		writeError(w, "invalid currency pair, expected format: USDRUB", http.StatusBadRequest)
		return
	}

	baseCode := strings.ToUpper(pair[:3])
	targetCode := strings.ToUpper(pair[3:])

	rateStr := r.FormValue("rate")
	if rateStr == "" {
		writeError(w, "rate is required", http.StatusBadRequest)
		return
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		writeError(w, "invalid rate value", http.StatusBadRequest)
		return
	}

	slog.InfoContext(ctx, "update exchange rate request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
		"base", baseCode,
		"target", targetCode,
		"rate", rate,
	)

	er, err := h.svc.Update(ctx, baseCode, targetCode, rate)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, "exchange rate not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(ctx, "failed to update exchange rate",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "handler",
			"error", err,
		)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(er)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
