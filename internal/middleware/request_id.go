package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// ctxKey — типизированный ключ для context.Value. Не голая строка,
// чтобы исключить коллизии с ключами из других пакетов.
type ctxKey string

const requestIDKey ctxKey = "request_id"

// RequestID генерит/пробрасывает сквозной идентификатор запроса.
// Если входящий X-Request-ID уже есть (запрос от API-шлюза или
// вышестоящего сервиса) — переиспользуем, иначе генерим UUID.
// Кладём ID:
//   - в context — для вызовов вниз по цепочке и для обработчиков
//     внутри сервиса;
//   - в заголовок ответа — чтобы клиент/саппорт видели ID в ответе.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IDFromContext достаёт request_id из context. Возвращает "" если
// RequestID middleware не отработал — это и есть нормальное
// значение по умолчанию.
func IDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
