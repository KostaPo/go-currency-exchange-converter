package exchange

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

func (h *Handler) Exchange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	from := strings.ToUpper(r.URL.Query().Get("from"))
	to := strings.ToUpper(r.URL.Query().Get("to"))
	amountStr := r.URL.Query().Get("amount")

	if from == "" || to == "" || amountStr == "" {
		writeError(w, "from, to and amount are required", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		writeError(w, "invalid amount value", http.StatusBadRequest)
		return
	}

	slog.InfoContext(ctx, "exchange request",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "handler",
		"from", from,
		"to", to,
		"amount", amount,
	)

	result, err := h.svc.Convert(ctx, from, to, amount)
	if err != nil {
		if errors.Is(err, ErrConversionNotPossible) {
			writeError(w, "exchange rate not found for this pair", http.StatusNotFound)
			return
		}
		slog.ErrorContext(ctx, "failed to convert currency",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "handler",
			"from", from,
			"to", to,
			"error", err,
		)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}
