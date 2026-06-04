package http

import (
	"net/http"

	"currency-exchange-converter/internal/currency"
	"currency-exchange-converter/internal/health"
)

func NewRouter(
	health *health.Handler,
	currencyHandler *currency.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	health.RegisterRoutes(mux, "/api/v1")

	currencyHandler.RegisterRoutes(mux, "/api/v1")

	return mux
}
