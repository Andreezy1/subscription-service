package repository

import (
	"errors"
	"subscription-service/internal/model"

	"github.com/jackc/pgx/v5"
)

func mapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrNotFound
	}
	return err
}
