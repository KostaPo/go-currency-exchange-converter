package health

import (
	"net/http"
)

func RegisterHealthRoutes(mux *http.ServeMux) {
	handler := NewHandler()

	mux.HandleFunc("GET /api/v1/health", handler.ServeHTTP)
}
