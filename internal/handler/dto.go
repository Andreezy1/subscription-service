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
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

func NewSubscriptionResponse(m *model.Subscription) *SubscriptionResponse {
	var endDatePtr *string

	if m.EndDate != nil {
		endDateStr := m.EndDate.Format("01-2006")
		endDatePtr = &endDateStr
	}
	return &SubscriptionResponse{
		ID:          m.ID,
		ServiceName: m.ServiceName,
		Price:       m.Price,
		UserID:      m.UserID,
		StartDate:   m.StartDate.Format("01-2006"),
		EndDate:     endDatePtr,
	}
}

func NewSubscriptionsResponse(subs []model.Subscription) []SubscriptionResponse {
	response := make([]SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		response = append(response, *NewSubscriptionResponse(&sub))
	}
	return response
}

type CalculateTotalCostResponse struct {
	Total int `json:"total"`
}

func NewCalculateTotalCostResponse(t int) *CalculateTotalCostResponse {
	return &CalculateTotalCostResponse{
		Total: t,
	}
}
