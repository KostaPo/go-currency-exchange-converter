# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s" \
    -o server \
    ./cmd/api

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata sqlite

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/server .
COPY --chown=app:app configs/ ./configs/
COPY --chown=app:app migrations/ ./migrations/
COPY --chown=app:app data/app.db ./data/app.db.seed
COPY --chown=app:app frontend/ ./frontend/

RUN chown -R app:app /app

USER app

EXPOSE 8080

VOLUME ["/app/data"]

ENTRYPOINT ["./server", "-config", "configs/config.docker.yml"]