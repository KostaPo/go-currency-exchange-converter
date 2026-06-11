package api

import (
	"net/http"

	"currency-exchange-converter/internal/middleware"
	"currency-exchange-converter/internal/model/currency"
	"currency-exchange-converter/internal/model/exchange"
	"currency-exchange-converter/internal/model/exchangerate"
	"currency-exchange-converter/internal/observability/health"
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
