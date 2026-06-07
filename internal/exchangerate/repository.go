package exchangerate

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"database/sql"
	"errors"
	"log/slog"
)

const (
	queryGetAll = `
	SELECT
		er.ID,
		er.Rate,

		bc.ID       AS base_id,
		bc.Code     AS base_code,
		bc.FullName AS base_full_name,
		bc.Sign     AS base_sign,

		tc.ID       AS target_id,
		tc.Code     AS target_code,
		tc.FullName AS target_full_name,
		tc.Sign     AS target_sign

	FROM ExchangeRates er
		JOIN Currencies bc ON bc.ID = er.BaseCurrencyId
		JOIN Currencies tc ON tc.ID = er.TargetCurrencyId

	ORDER BY er.ID
`
)

var (
	ErrNotFound      = errors.New("currency not found")
	ErrAlreadyExists = errors.New("currency already exists")

	ErrInvalidID       = errors.New("invalid currency id")
	ErrInvalidCode     = errors.New("currency code must be 3 uppercase letters")
	ErrInvalidFullName = errors.New("currency full name is required")
)

type Repository interface {
	GetAll(ctx context.Context) ([]*ExchangeRate, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]*ExchangeRate, error) {
	slog.DebugContext(ctx, "executing query",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "GetAll",
	)

	rows, err := r.db.QueryContext(ctx, queryGetAll)
	if err != nil {
		logRepoError(ctx, "GetAll", err)
		return nil, err
	}
	defer rows.Close()

	var result []*ExchangeRate

	for rows.Next() {
		var er ExchangeRate

		if err := rows.Scan(
			&er.ID,
			&er.Rate,
			&er.BaseCurrency.ID,
			&er.BaseCurrency.Code,
			&er.BaseCurrency.FullName,
			&er.BaseCurrency.Sign,
			&er.TargetCurrency.ID,
			&er.TargetCurrency.Code,
			&er.TargetCurrency.FullName,
			&er.TargetCurrency.Sign,
		); err != nil {
			logRepoError(ctx, "GetAll.Scan", err)
			return nil, err
		}

		result = append(result, &er)
	}

	if err := rows.Err(); err != nil {
		logRepoError(ctx, "GetAll.Rows", err)
		return nil, err
	}

	slog.DebugContext(ctx, "query done",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "GetAll",
		"rows", len(result),
	)

	return result, nil
}

// logRepoError — хелпер для единообразного логирования ошибок репозитория.
// Различает отмену клиентом (Warn) и реальную ошибку БД (Error) —
func logRepoError(ctx context.Context, query string, err error) {
	if errors.Is(err, context.Canceled) {
		slog.WarnContext(ctx, "client disconnected",
			"request_id", middleware.IDFromContext(ctx),
			"layer", "repository",
			"query", query,
		)
		return
	}
	slog.ErrorContext(ctx, "query failed",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", query,
		"error", err,
	)
}
