package api

import (
	"net/http"

	"currency-exchange-converter/internal/api/health"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	health.RegisterHealthRoutes(mux)

	return mux
}
