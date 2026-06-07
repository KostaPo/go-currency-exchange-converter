package exchangerate

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.HandleFunc("GET "+prefix+"/exchangeRates", h.GetAll)
	mux.HandleFunc("GET "+prefix+"/exchangeRate/{pair}", h.GetByPair)
}
