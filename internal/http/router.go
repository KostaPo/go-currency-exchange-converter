package http

import (
	"net/http"

	"currency-exchange-converter/internal/health"
)

type Handlers struct {
	Health *health.Handler
}

func NewRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /api/v1/health",
		h.Health.GetHealth,
	)

	return mux
}
