package exchange

import (
	"context"
	"currency-exchange-converter/internal/exchangerate"
	"currency-exchange-converter/internal/middleware"
	"errors"
	"log/slog"
)

var ErrConversionNotPossible = errors.New("conversion not possible")

type Service interface {
	Convert(ctx context.Context, from, to string, amount float64) (*Exchange, error)
}

type service struct {
	exchangeRateRepo exchangerate.Repository
}

func NewService(exchangeRateRepo exchangerate.Repository) Service {
	return &service{exchangeRateRepo: exchangeRateRepo}
}

func (s *service) Convert(ctx context.Context, from, to string, amount float64) (*Exchange, error) {
	slog.InfoContext(ctx, "converting currency",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"from", from,
		"to", to,
		"amount", amount,
	)

	// сценарий 1 — прямая пара AB
	er, err := s.exchangeRateRepo.GetByPair(ctx, from, to)
	if err == nil {
		slog.InfoContext(ctx, "direct pair found",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"from", from,
			"to", to,
		)
		return &Exchange{
			BaseCurrency:    er.BaseCurrency,
			TargetCurrency:  er.TargetCurrency,
			Rate:            er.Rate,
			Amount:          amount,
			ConvertedAmount: amount * er.Rate,
		}, nil
	}

	// сценарий 2 — обратная пара BA
	er, err = s.exchangeRateRepo.GetByPair(ctx, to, from)
	if err == nil {
		rate := 1 / er.Rate
		slog.InfoContext(ctx, "reverse pair found",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"from", from,
			"to", to,
			"rate", rate,
		)
		return &Exchange{
			BaseCurrency:    er.TargetCurrency, // меняем местами
			TargetCurrency:  er.BaseCurrency,
			Rate:            rate,
			Amount:          amount,
			ConvertedAmount: amount * rate,
		}, nil
	}

	// сценарий 3 — кросс-курс через USD
	usdToFrom, err := s.exchangeRateRepo.GetByPair(ctx, "USD", from)
	if err != nil {
		slog.InfoContext(ctx, "USD pair not found",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"pair", "USD/"+from,
		)
		return nil, ErrConversionNotPossible
	}

	usdToTarget, err := s.exchangeRateRepo.GetByPair(ctx, "USD", to)
	if err != nil {
		slog.InfoContext(ctx, "USD pair not found",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"pair", "USD/"+to,
		)
		return nil, ErrConversionNotPossible
	}

	// rateAB = rateUSD_B / rateUSD_A
	rate := usdToTarget.Rate / usdToFrom.Rate
	slog.InfoContext(ctx, "cross rate calculated",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"from", from,
		"to", to,
		"rate", rate,
	)

	return &Exchange{
		BaseCurrency:    usdToFrom.TargetCurrency,
		TargetCurrency:  usdToTarget.TargetCurrency,
		Rate:            rate,
		Amount:          amount,
		ConvertedAmount: amount * rate,
	}, nil
}
