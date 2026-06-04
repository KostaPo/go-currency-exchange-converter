package currency

import (
	"context"
	"database/sql"
	"errors"
)

const (
	queryGetAll = `
        SELECT ID, Code, FullName, Sign
		FROM Currencies
		ORDER BY Code
    `

	queryGetByID = `
        SELECT ID, Code, FullName, Sign
        FROM currencies
        WHERE id = $1
    `

	queryGetByCode = `
        SELECT ID, Code, FullName, Sign
        FROM currencies
        WHERE code = $1
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
	GetByID(ctx context.Context, id int) (*Currency, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]*Currency, error) {

	rows, err := r.db.QueryContext(ctx, queryGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Currency

	for rows.Next() {
		var c Currency

		if err := rows.Scan(
			&c.ID,
			&c.Code,
			&c.FullName,
			&c.Sign,
		); err != nil {
			return nil, err
		}

		result = append(result, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) Create(ctx context.Context, c *Currency) error {
	// реализация запроса к БД
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*Currency, error) {
	// реализация запроса к БД
	return nil, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	// реализация запроса к БД
	return nil
}
