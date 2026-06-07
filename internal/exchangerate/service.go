package exchangerate

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"errors"
	"log/slog"
)

type Service interface {
	GetAll(ctx context.Context) ([]*ExchangeRate, error)
	GetByPair(ctx context.Context, baseCode, targetCode string) (*ExchangeRate, error)
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

func (s *service) GetByPair(ctx context.Context, baseCode, targetCode string) (*ExchangeRate, error) {
	slog.InfoContext(ctx, "fetching exchange rate by pair",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"base", baseCode,
		"target", targetCode,
	)

	rate, err := s.repo.GetByPair(ctx, baseCode, targetCode)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			slog.InfoContext(ctx, "exchange rate not found",
				"request_id", middleware.IDFromContext(ctx),
				"layer", "service",
				"base", baseCode,
				"target", targetCode,
			)
			return nil, ErrNotFound
		}
		slog.ErrorContext(ctx, "failed to fetch exchange rate",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"base", baseCode,
			"target", targetCode,
			"error", err,
		)
		return nil, err
	}

	slog.InfoContext(ctx, "exchange rate fetched",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"base", baseCode,
		"target", targetCode,
	)

	return rate, nil
}
