package model

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

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

func TestSubscription_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sub     Subscription
		wantErr bool
	}{
		{
			name: "valid subscription",
			sub: Subscription{
				ServiceName: "Netflix",
				Price:       100,
				StartDate:   date(2025, time.January),
				EndDate:     ptr(date(2025, time.December)),
			},
			wantErr: false,
		},
		{
			name: "valid without end date",
			sub: Subscription{
				ServiceName: "Netflix",
				Price:       100,
				StartDate:   date(2025, time.January),
				EndDate:     nil,
			},
			wantErr: false,
		},
		{
			name: "empty service name",
			sub: Subscription{
				ServiceName: "",
				Price:       100,
				StartDate:   date(2025, time.January),
			},
			wantErr: true,
		},
		{
			name: "zero price",
			sub: Subscription{
				ServiceName: "Netflix",
				Price:       0,
				StartDate:   date(2025, time.January),
			},
			wantErr: true,
		},
		{
			name: "negative price",
			sub: Subscription{
				ServiceName: "Netflix",
				Price:       -100,
				StartDate:   date(2025, time.January),
			},
			wantErr: true,
		},
		{
			name: "start date after end date",
			sub: Subscription{
				ServiceName: "Netflix",
				Price:       100,
				StartDate:   date(2025, time.December),
				EndDate:     ptr(date(2025, time.January)),
			},
			wantErr: true,
		},
		{
			name: "nil uuid",
			sub: Subscription{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.Nil,
				StartDate:   date(2025, time.January),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sub.Validate()

			if tt.wantErr && !errors.Is(err, ErrValidate) {
				t.Fatalf("expected ErrValidate, got %v", err)
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	}
}
