package handler

import (
	"encoding/json"
	"net/http"
)

func NewRouter(
// exchangeRate *ExchangeRateHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	registerHealthRoutes(mux)

	// exchangeRate.RegisterRoutes(
	// 	mux,
	// 	"/api/v1/exchange-rates",
	// )

	return mux
}

func registerHealthRoutes(
	mux *http.ServeMux,
) {

	mux.HandleFunc(
		"GET /health",
		healthHandler,
	)
}

func healthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(
		map[string]string{
			"status": "ok",
		},
	)
}
