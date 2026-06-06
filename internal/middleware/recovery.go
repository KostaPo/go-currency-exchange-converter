package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery ловит панику в любом обработчике ниже по цепочке,
// логирует то чем паниковали + стек, и отдаёт клиенту 500.
// Канонически ставится САМЫМ ВНЕШНИМ middleware — иначе паника
// в Logger/RequestID его обойдёт.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					"panic", rec,
					"method", r.Method,
					"path", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
