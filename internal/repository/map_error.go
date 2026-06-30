package repository

import (
	"errors"
	"fmt"
	"subscription-service/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func mapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation: // 23505 → 409
			return fmt.Errorf("%w: %s", model.ErrConflict, pgErr.Message)
		case pgerrcode.CheckViolation: // 23514 → 400
			return fmt.Errorf("%w: %s", model.ErrValidate, pgErr.Message)
		}
	}
	return err
}
