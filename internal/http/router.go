package http

import (
	"net/http"

	"currency-exchange-converter/internal/health"
)

func NewRouter(
	health *health.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	health.RegisterRoutes(mux, "/api/v1")

	return mux
}
