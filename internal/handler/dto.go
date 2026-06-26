package handler

import (
	"subscription-service/internal/model"
	"time"

	"github.com/google/uuid"
)

type SubscriptionRequest struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

func (r *SubscriptionRequest) ToModel() (*model.Subscription, error) {
	startDate, err := parseMonth(r.StartDate)
	if err != nil {
		return nil, err
	}
	var endDate *time.Time
	if r.EndDate != nil {
		parsedEndDate, err := parseMonth(*r.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &parsedEndDate
	}
	return &model.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      r.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

type SubscriptionResponse struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

func NewSubscriptionResponse(m *model.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:          m.ID,
		ServiceName: m.ServiceName,
		Price:       m.Price,
		UserID:      m.UserID,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
	}
}

type CalculateTotalCostResponse struct {
	Total int `json:"total"`
}

func NewCalculateTotalCostResponse(t int) *CalculateTotalCostResponse {
	return &CalculateTotalCostResponse{
		Total: t,
	}
}
