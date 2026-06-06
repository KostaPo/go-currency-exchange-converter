package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter оборачивает http.ResponseWriter, чтобы перехватить
// статус-код для логов. Канонический обходной приём: интерфейс
// http.ResponseWriter не предоставляет метод Status(), а реальная
// реализация под капотом неэкспортируемая, и приведение типа не
// сработает. Поэтому ловим статус через переопределённый
// WriteHeader.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger пишет структурированный лог по каждому запросу: method,
// path, status, duration_ms. Использует stdlib log/slog — без
// zap/zerolog. Статус по умолчанию http.StatusOK: если обработчик
// не вызвал WriteHeader, Go всё равно вернёт 200, и мы должны это
// отразить в логе.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		// request_id берём из r.Context(): RequestID middleware стоит
		// СНАРУЖИ Logger в цепочке, поэтому к этому моменту он уже
		// положил ID в context того r, который сюда дошёл.
		slog.Info("http",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", IDFromContext(r.Context()),
		)
	})
}
