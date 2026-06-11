package currency

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.HandleFunc("GET "+prefix+"/currencies", h.GetAll)
	mux.HandleFunc("POST "+prefix+"/currencies", h.Create)
	mux.HandleFunc("GET "+prefix+"/currency/{code}", h.GetByCode)
	mux.HandleFunc("GET "+prefix+"/currency/", h.handleEmptyCode)
}
