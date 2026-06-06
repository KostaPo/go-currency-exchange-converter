package currency

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"log/slog"
)

type Service interface {
	GetAll(ctx context.Context) ([]*Currency, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]*Currency, error) {

	// Service логирует бизнес-событие — не HTTP-детали, а то что
	// происходит с данными. request_id тот же что в handler —
	// контекст прокинут сквозь все слои.
	slog.InfoContext(ctx, "fetching all currencies",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
	)

	currencies, err := s.repo.GetAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch currencies",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"error", err,
		)
		return nil, err
	}

	slog.InfoContext(ctx, "currencies fetched",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"count", len(currencies),
	)

	return currencies, nil
}
