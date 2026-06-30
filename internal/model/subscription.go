package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int        `db:"price"`
	UserID      uuid.UUID  `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}

func (s *Subscription) Validate() error {
	if s.ServiceName == "" {
		return fmt.Errorf("%w: service_name is empty", ErrValidate)
	}
	if s.Price <= 0 {
		return fmt.Errorf("%w: price must be positive", ErrValidate)
	}
	if s.UserID == uuid.Nil {
		return fmt.Errorf("%w: user_id is required", ErrValidate)
	}
	if s.EndDate != nil {
		if s.StartDate.After(*s.EndDate) {
			return fmt.Errorf("%w: start_date is after end_date", ErrValidate)
		}
	}
	return nil
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
}

type CalculateTotalCostFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartDate   time.Time
	EndDate     time.Time
}
