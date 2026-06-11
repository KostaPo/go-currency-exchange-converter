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
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
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
		writeError(w, "failed to get currencies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(currencies)
}

func (h *Handler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	code := r.PathValue("code")
	if code == "" {
		writeError(w, "code is required", http.StatusBadRequest)
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
			writeError(w, "currency not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(ctx, "failed to get currency",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "handler",
			"code", code,
			"error", err,
		)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(currency)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := r.ParseForm(); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	fullName := r.FormValue("name")
	sign := r.FormValue("sign")

	slog.InfoContext(ctx, "create currency request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
		"code", code,
	)

	c, err := h.svc.Create(ctx, code, fullName, sign)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidSign):
			writeError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrInvalidCode),
			errors.Is(err, ErrInvalidFullName):
			writeError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrAlreadyExists):
			writeError(w, "currency already exists", http.StatusConflict)
		default:
			writeError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

func (h *Handler) handleEmptyCode(w http.ResponseWriter, r *http.Request) {
	writeError(w, "currency code is required", http.StatusBadRequest)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
