package api

import (
	"net/http"

	"currency-exchange-converter/internal/currency"
	"currency-exchange-converter/internal/exchange"
	"currency-exchange-converter/internal/exchangerate"
	"currency-exchange-converter/internal/health"
	"currency-exchange-converter/internal/middleware"
)

func NewRouter(
	health *health.Handler,
	currencyHandler *currency.Handler,
	exchangerateHandler *exchangerate.Handler,
	exchangeHandler *exchange.Handler,
) http.Handler {
	mux := http.NewServeMux()

	health.RegisterRoutes(mux, "/api/v1")
	currencyHandler.RegisterRoutes(mux, "/api/v1")
	exchangerateHandler.RegisterRoutes(mux, "/api/v1")
	exchangeHandler.RegisterRoutes(mux, "/api/v1")

	// Раздаём фронтенд
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	return middleware.Recovery(
		middleware.RequestID(
			middleware.Logger(mux)))

}
