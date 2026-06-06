package http

import (
	"net/http"

	"currency-exchange-converter/internal/currency"
	"currency-exchange-converter/internal/health"
	"currency-exchange-converter/internal/middleware"
)

func NewRouter(
	health *health.Handler,
	currencyHandler *currency.Handler,
) http.Handler {
	mux := http.NewServeMux()

	health.RegisterRoutes(mux, "/api/v1")

	currencyHandler.RegisterRoutes(mux, "/api/v1")

	return middleware.Recovery(
		middleware.RequestID(
			middleware.Logger(mux)))

}
