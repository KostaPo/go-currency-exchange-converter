package exchangerate

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
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
	queryGetByPair = `
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

	WHERE bc.Code = ? AND tc.Code = ?
	`

	queryCreate = `
	INSERT INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate)
	VALUES (
		(SELECT ID FROM Currencies WHERE Code = ?),
		(SELECT ID FROM Currencies WHERE Code = ?),
		?
	)
	`

	queryUpdate = `
	UPDATE ExchangeRates
	SET Rate = ?
	WHERE BaseCurrencyId = (SELECT ID FROM Currencies WHERE Code = ?)
	AND TargetCurrencyId = (SELECT ID FROM Currencies WHERE Code = ?)
`
)

var (
	ErrNotFound         = errors.New("exchange rate not found")
	ErrAlreadyExists    = errors.New("ExchangeRate already exists")
	ErrCurrencyNotFound = errors.New("currency not found")
)

type Repository interface {
	GetAll(ctx context.Context) ([]*ExchangeRate, error)
	GetByPair(ctx context.Context, baseCode, targetCode string) (*ExchangeRate, error)
	Create(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error)
	Update(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]*ExchangeRate, error) {
	slog.DebugContext(ctx, "executing...",
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

func (r *repository) GetByPair(ctx context.Context, baseCode, targetCode string) (*ExchangeRate, error) {
	slog.DebugContext(ctx, "executing...",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "GetByPair",
		"base", baseCode,
		"target", targetCode,
	)

	var er ExchangeRate
	err := r.db.QueryRowContext(ctx, queryGetByPair, baseCode, targetCode).Scan(
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
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		logRepoError(ctx, "GetByPair", err)
		return nil, err
	}

	slog.DebugContext(ctx, "query done",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "GetByPair",
		"base", baseCode,
		"target", targetCode,
	)

	return &er, nil

}

func (r *repository) Create(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error) {
	slog.DebugContext(ctx, "executing...",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "Create",
	)

	_, err := r.db.ExecContext(ctx, queryCreate, baseCode, targetCode, rate)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, ErrAlreadyExists
		}
		logRepoError(ctx, "Create.Insert", err)
		return nil, err
	}

	return r.GetByPair(ctx, baseCode, targetCode)
}

func (r *repository) Update(ctx context.Context, baseCode, targetCode string, rate float64) (*ExchangeRate, error) {
	slog.DebugContext(ctx, "executing...",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "Update",
	)

	result, err := r.db.ExecContext(ctx, queryUpdate, rate, baseCode, targetCode)
	if err != nil {
		logRepoError(ctx, "Update", err)
		return nil, err
	}
	if err != nil {
		logRepoError(ctx, "Update", err)
		return nil, err
	}

	// если RowsAffected == 0 — пара не найдена
	rows, err := result.RowsAffected()
	if err != nil {
		logRepoError(ctx, "Update.RowsAffected", err)
		return nil, err
	}
	if rows == 0 {
		return nil, ErrNotFound
	}

	return r.GetByPair(ctx, baseCode, targetCode)
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
