package currency

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"errors"
	"log/slog"
	"regexp"
	"strings"
)

var codeChecker = regexp.MustCompile(`^[A-Z]{3}$`)

type Service interface {
	GetAll(ctx context.Context) ([]*Currency, error)
	GetByCode(ctx context.Context, code string) (*Currency, error)
	Create(ctx context.Context, code, fullName, sign string) (*Currency, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]*Currency, error) {

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

func (s *service) GetByCode(ctx context.Context, code string) (*Currency, error) {

	slog.InfoContext(ctx, "fetching currency by code",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
	)

	currency, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			slog.InfoContext(ctx, "currency not found",
				"request_id", middleware.IDFromContext(ctx),
				"layer", "service",
				"code", code,
			)
			return nil, ErrNotFound
		}
		slog.ErrorContext(ctx, "failed to fetch currency",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "service",
			"code", code,
			"error", err,
		)
		return nil, err
	}

	slog.InfoContext(ctx, "currency fetched",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"code", code,
	)

	return currency, nil
}

func (s *service) Create(ctx context.Context, code, name, sign string) (*Currency, error) {
	if !codeChecker.MatchString(code) {
		return nil, ErrInvalidCode
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidFullName
	}
	if strings.TrimSpace(sign) == "" { // <- TC-024, TC-025
		return nil, ErrInvalidSign
	}
	if len([]rune(sign)) > 3 { // <- TC-025a
		return nil, ErrInvalidSign
	}

	slog.InfoContext(ctx, "creating currency",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"code", code,
	)

	c := &Currency{Code: code, FullName: name, Sign: sign}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "currency created",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "service",
		"code", code,
	)

	return c, nil
}
