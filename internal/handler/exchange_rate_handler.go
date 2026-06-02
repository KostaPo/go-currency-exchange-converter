package handler

import "net/http"

type ExchangeRateHandler struct{}

func NewExchangeRateHandler() *ExchangeRateHandler {
	return &ExchangeRateHandler{}
}

func (h *ExchangeRateHandler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.HandleFunc(prefix, h.list)
	mux.HandleFunc("POST "+prefix, h.create)
	mux.HandleFunc(prefix+"/{id}", h.get)
	mux.HandleFunc("PUT "+prefix+"/{id}", h.update)
	mux.HandleFunc("DELETE "+prefix+"/{id}", h.delete)
}

func (h *ExchangeRateHandler) list(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list placeholder"))
}

func (h *ExchangeRateHandler) create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create placeholder"))
}

func (h *ExchangeRateHandler) get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get placeholder"))
}

func (h *ExchangeRateHandler) update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update placeholder"))
}

func (h *ExchangeRateHandler) delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete placeholder"))
}
