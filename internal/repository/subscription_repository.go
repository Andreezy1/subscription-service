package repository

import (
	"context"
	"fmt"
	"strings"
	"subscription-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) Create(ctx context.Context, sub *model.Subscription) error {
	query := `
			INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			VALUES (@service_name, @price, @user_id, @start_date, @end_date)
			RETURNING id;
			`
	args := pgx.NamedArgs{
		"service_name": sub.ServiceName,
		"price":        sub.Price,
		"user_id":      sub.UserID,
		"start_date":   sub.StartDate,
		"end_date":     sub.EndDate,
	}
	row := r.db.QueryRow(ctx, query, args)

	if err := row.Scan(&sub.ID); err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*model.Subscription, error) {
	query := `
			SELECT 
				id,
				service_name,
				price,
				user_id,
				start_date,
				end_date
			FROM subscriptions
			WHERE id = $1;
			`
	row := r.db.QueryRow(ctx, query, id)
	sub := &model.Subscription{}
	err := row.Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		return nil, mapError(err)
	}
	return sub, nil
}

func (r *postgresRepository) Update(ctx context.Context, sub *model.Subscription) error {
	query := `
			UPDATE subscriptions
			SET 
				service_name = $1,
				price = $2,
				user_id = $3,
				start_date = $4,
				end_date = $5
			WHERE id = $6;
			`
	cmdTag, err := r.db.Exec(ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return mapError(err)
	}
	if cmdTag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `
			DELETE FROM subscriptions
			WHERE id = $1;
			`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return mapError(err)
	}
	if cmdTag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *postgresRepository) List(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error) {
	query := `
			SELECT 
				id,
				service_name,
				price,
				user_id,
				start_date,
				end_date
			FROM subscriptions
			`
	conditions := []string{}
	args := []any{}
	if filter.UserID != nil {
		conditions = append(
			conditions,
			fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, *filter.UserID)
	}
	if filter.ServiceName != nil {
		conditions = append(
			conditions,
			fmt.Sprintf("service_name = $%d", len(args)+1))
		args = append(args, *filter.ServiceName)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	subscriptions := make([]model.Subscription, 0)
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, mapError(err)
		}
		subscriptions = append(subscriptions, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return subscriptions, nil
}
