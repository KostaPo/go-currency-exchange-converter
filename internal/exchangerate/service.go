package exchangerate

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"log/slog"
)

type Service interface {
	GetAll(ctx context.Context) ([]*ExchangeRate, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]*ExchangeRate, error) {
	slog.InfoContext(ctx, "fetching all exchange rates",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
	)

	rates, err := s.repo.GetAll(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch exchange rates",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"error", err,
		)
		return nil, err
	}

	slog.InfoContext(ctx, "exchange rates fetched",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"count", len(rates),
	)

	return rates, nil
}
