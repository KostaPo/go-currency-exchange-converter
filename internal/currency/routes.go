package currency

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.HandleFunc("GET "+prefix+"/currencies", h.GetAll)
}
