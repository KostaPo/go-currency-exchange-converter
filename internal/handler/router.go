package handler

import "net/http"

func NewRouter(
	exchangeRate *ExchangeRateHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	exchangeRate.RegisterRoutes(mux, "/api/v1/exchange-rates")

	return mux
}
