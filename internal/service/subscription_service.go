package service

import (
	"context"
	"fmt"
	"subscription-service/internal/model"
	"time"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id int) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error)
}

type subscriptionService struct {
	repo SubscriptionRepository
}

func NewSubscriptionService(repo SubscriptionRepository) *subscriptionService {
	return &subscriptionService{
		repo: repo,
	}
}

func (s *subscriptionService) Create(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	if err := sub.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id int) (*model.Subscription, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id < 0", model.ErrValidate)
	}
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *subscriptionService) Update(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	if err := sub.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *subscriptionService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: id<0", model.ErrValidate)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *subscriptionService) List(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error) {
	return s.repo.List(ctx, filter)
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, filter model.CalculateTotalCostFilter) (int, error) {
	if filter.StartDate.After(filter.EndDate) {
		return 0, fmt.Errorf("%w: start date after end date", model.ErrValidate)
	}
	subs, err := s.repo.List(ctx, model.SubscriptionFilter{
		UserID:      filter.UserID,
		ServiceName: filter.ServiceName,
	})
	if err != nil {
		return 0, err
	}
	total := 0
	for _, sub := range subs {
		months := calculateMonths(
			&sub,
			filter.StartDate,
			filter.EndDate,
		)
		total += months * sub.Price
	}
	return total, nil
}

func calculateMonths(sub *model.Subscription, from time.Time, to time.Time) int {
	subEnd := to
	if sub.EndDate != nil {
		subEnd = *sub.EndDate
	}

	start := maxTime(from, sub.StartDate)
	end := minTime(to, subEnd)

	if start.After(end) {
		return 0
	}
	return ((end.Year()-start.Year())*12 + (int(end.Month()) - int(start.Month())) + 1)
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
