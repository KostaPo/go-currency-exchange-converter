package exchangerate

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"currency-exchange-converter/internal/model/currency"
	"errors"
	"fmt"
	"log/slog"
)

type Service interface {
	GetAll(ctx context.Context) ([]*ExchangeRate, error)
	GetByPair(ctx context.Context, baseCode, targetCode string) (*ExchangeRate, error)
	Create(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error)
	Update(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error)
}

type service struct {
	repo         Repository
	currencyRepo currency.Repository
}

func NewService(repo Repository, currencyRepo currency.Repository) Service {
	return &service{repo: repo, currencyRepo: currencyRepo}
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
	)

	return rate, nil
}

func (s *service) Create(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error) {
	slog.InfoContext(ctx, "creating exchange rate",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
	)

	// проверяем базовую валюту
	_, err := s.currencyRepo.GetByCode(ctx, baseCode)
	if err != nil {
		if errors.Is(err, currency.ErrNotFound) {
			return nil, fmt.Errorf("base currency %s not found", baseCode)
		}
		return nil, err
	}

	// проверяем целевую валюту
	_, err = s.currencyRepo.GetByCode(ctx, targetCode)
	if err != nil {
		if errors.Is(err, currency.ErrNotFound) {
			return nil, fmt.Errorf("target currency %s not found", targetCode)
		}
		return nil, err
	}

	er, err := s.repo.Create(ctx, baseCode, targetCode, rate)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return nil, ErrAlreadyExists
		}
		slog.ErrorContext(ctx, "failed to create exchange rate",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"error", err,
		)
		return nil, err
	}

	slog.InfoContext(ctx, "exchange rate created",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
	)

	return er, nil
}

func (s *service) Update(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error) {
	slog.InfoContext(ctx, "updating exchange rate",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
	)

	er, err := s.repo.Update(ctx, baseCode, targetCode, rate)
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
		slog.ErrorContext(ctx, "failed to update exchange rate",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"error", err,
		)
		return nil, err
	}

	slog.InfoContext(ctx, "exchange rate updated",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"base", baseCode,
		"target", targetCode,
	)

	return er, nil
}
