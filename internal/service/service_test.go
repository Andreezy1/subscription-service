package service

import (
	"context"
	"errors"
	"subscription-service/internal/model"
	"testing"
	"time"
)

type repoMock struct {
	listFn func(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error)
}

func (m *repoMock) List(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error) {
	return m.listFn(ctx, filter)
}

func (m *repoMock) Create(ctx context.Context, sub *model.Subscription) error {
	panic("not implemented")
}

func (m *repoMock) GetByID(ctx context.Context, id int) (*model.Subscription, error) {
	panic("not implemented")
}

func (m *repoMock) Update(ctx context.Context, sub *model.Subscription) error {
	panic("not implemented")
}

func (m *repoMock) Delete(ctx context.Context, id int) error {
	panic("not implemented")
}

func date(year int, month time.Month) time.Time {
	return time.Date(
		year,
		month,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}
func ptr(t time.Time) *time.Time {
	return &t
}

func TestCalculateMonths(t *testing.T) {
	tests := []struct {
		name string

		sub  model.Subscription
		from time.Time
		to   time.Time

		want int
	}{
		{
			name: "full intersection",
			sub: model.Subscription{
				StartDate: date(2025, time.July),
				EndDate:   ptr(date(2025, time.September)),
			},
			from: date(2025, time.January),
			to:   date(2025, time.December),
			want: 3,
		},
		{
			name: "single month overlap",
			sub: model.Subscription{
				StartDate: date(2025, time.July),
				EndDate:   ptr(date(2025, time.September)),
			},
			from: date(2025, time.August),
			to:   date(2025, time.August),
			want: 1,
		},
		{
			name: "no overlap",
			sub: model.Subscription{
				StartDate: date(2025, time.July),
				EndDate:   ptr(date(2025, time.September)),
			},
			from: date(2026, time.January),
			to:   date(2026, time.March),
			want: 0,
		},
		{
			name: "full overlap",
			sub: model.Subscription{
				StartDate: date(2025, time.July),
				EndDate:   nil,
			},
			from: date(2026, time.January),
			to:   date(2026, time.March),
			want: 3,
		},
		{
			name: "partial overlap",
			sub: model.Subscription{
				StartDate: date(2025, time.July),
				EndDate:   nil,
			},
			from: date(2025, time.June),
			to:   date(2025, time.August),
			want: 2,
		},
		{
			name: "request inside subscription",
			sub: model.Subscription{
				StartDate: date(2025, time.July),
				EndDate:   ptr(date(2025, time.December)),
			},
			from: date(2025, time.September),
			to:   date(2025, time.October),
			want: 2,
		},
		{
			name: "subscription boundaries",
			sub: model.Subscription{
				StartDate: date(2025, time.July),
				EndDate:   ptr(date(2025, time.July)),
			},
			from: date(2025, time.July),
			to:   date(2025, time.July),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateMonths(
				&tt.sub,
				tt.from,
				tt.to,
			)

			if got != tt.want {
				t.Errorf(
					"calculateMonths(%s) = %d, want %d",
					tt.name,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCalculateTotalCost(t *testing.T) {
	tests := []struct {
		name    string
		filter  model.CalculateTotalCostFilter
		subs    []model.Subscription
		repoErr error

		wantTotal int
		wantErr   bool
	}{
		{
			name: "Single subscription",
			filter: model.CalculateTotalCostFilter{
				StartDate: date(2025, time.January),
				EndDate:   date(2025, time.March),
			},
			subs: []model.Subscription{
				{
					ServiceName: "Netflix",
					Price:       100,
					StartDate:   date(2025, time.January),
					EndDate:     ptr(date(2025, time.March)),
				},
			}, repoErr: nil,
			wantTotal: 300,
			wantErr:   false,
		},
		{
			name: "Multiple subscriptions",
			filter: model.CalculateTotalCostFilter{
				StartDate: date(2025, time.January),
				EndDate:   date(2025, time.March),
			},
			subs: []model.Subscription{
				{
					ServiceName: "Netflix",
					Price:       100,
					StartDate:   date(2025, time.January),
					EndDate:     ptr(date(2025, time.March)),
				},
				{
					ServiceName: "Spotify",
					Price:       200,
					StartDate:   date(2025, time.January),
					EndDate:     ptr(date(2025, time.March)),
				},
			}, repoErr: nil,
			wantTotal: 900,
			wantErr:   false,
		},
		{
			name: "Empty list",
			filter: model.CalculateTotalCostFilter{
				StartDate: date(2025, time.January),
				EndDate:   date(2025, time.March),
			},
			subs: []model.Subscription{
				{},
			}, repoErr: nil,
			wantTotal: 0,
			wantErr:   false,
		},
		{
			name: "Repository error",
			filter: model.CalculateTotalCostFilter{
				StartDate: date(2025, time.January),
				EndDate:   date(2025, time.March),
			},
			subs:      nil,
			repoErr:   errors.New("db error"),
			wantTotal: 0,
			wantErr:   true,
		},
		{
			name: "Invalid period",
			filter: model.CalculateTotalCostFilter{
				StartDate: date(2025, time.December),
				EndDate:   date(2025, time.March),
			},
			subs: []model.Subscription{
				{
					ServiceName: "Netflix",
					Price:       100,
					StartDate:   date(2025, time.January),
					EndDate:     ptr(date(2025, time.March)),
				},
			},
			repoErr:   nil,
			wantTotal: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoMock{
				listFn: func(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error) {
					return tt.subs, tt.repoErr
				},
			}
			service := NewSubscriptionService(repo)
			total, err := service.CalculateTotalCost(context.Background(), tt.filter)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if total != tt.wantTotal {
				t.Fatalf(
					"expected total %d, got %d",
					tt.wantTotal,
					total)
			}
			if tt.name == "invalid period" {
				if !errors.Is(err, model.ErrValidate) {
					t.Fatalf("expected ErrValidate, got %v", err)
				}
				return
			}
		})
	}
}
