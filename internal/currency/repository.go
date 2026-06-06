package currency

import (
	"context"
	"currency-exchange-converter/internal/middleware"
	"database/sql"
	"errors"
	"log/slog"
)

const (
	queryGetAll = `
        SELECT ID, Code, FullName, Sign
		FROM Currencies
		ORDER BY Code
    `

	queryGetByCode = `
        SELECT ID, Code, FullName, Sign
    	FROM Currencies
    	WHERE Code = ?
    `

	queryCreate = `
        INSERT INTO currencies (Code, FullName, Sign)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	queryDelete = `
        DELETE FROM currencies
        WHERE id = $1
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
	Create(ctx context.Context, c *Currency) error
	GetAll(ctx context.Context) ([]*Currency, error)
	GetByCode(ctx context.Context, code string) (*Currency, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]*Currency, error) {

	// Repository логирует технические детали — SQL-операции.
	// Debug уровень: в продакшне молчит, включается только при отладке.
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

	var result []*Currency

	for rows.Next() {
		var c Currency
		if err := rows.Scan(&c.ID, &c.Code, &c.FullName, &c.Sign); err != nil {
			slog.ErrorContext(ctx, "failed to scan row",
				"request_id", middleware.IDFromContext(ctx),
				"layer", "repository",
				"error", err,
			)
			return nil, err
		}
		result = append(result, &c)
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

func (r *repository) GetByCode(ctx context.Context, code string) (*Currency, error) {
	slog.DebugContext(ctx, "executing query",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "GetByCode",
		"code", code,
	)

	var c Currency
	err := r.db.QueryRowContext(ctx, queryGetByCode, code).Scan(
		&c.ID,
		&c.Code,
		&c.FullName,
		&c.Sign,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		logRepoError(ctx, "GetByCode", err)
		return nil, err
	}

	slog.DebugContext(ctx, "query done",
		"request_id", middleware.IDFromContext(ctx),
		"layer", "repository",
		"query", "GetByCode",
		"code", code,
	)

	return &c, nil
}

func (r *repository) Create(ctx context.Context, c *Currency) error {
	// реализация запроса к БД
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	// реализация запроса к БД
	return nil
}

// --- helper ---

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
