package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"subscription-service/internal/model"
	"testing"
	"time"
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

func TestParseQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name: "valid request",
			url: "/subscriptions/total?" +
				"start_date=01-2025&" +
				"end_date=12-2025&" +
				"user_id=550e8400-e29b-41d4-a716-446655440000&" +
				"service_name=Netflix",
			wantErr: false,
		},
		{
			name: "invalid start date",
			url: "/subscriptions/total?" +
				"start_date=abc&" +
				"end_date=12-2025",
			wantErr: true,
		},
		{
			name: "invalid end date",
			url: "/subscriptions/total?" +
				"start_date=01-2025&" +
				"end_date=xyz",
			wantErr: true,
		},
		{
			name: "invalid uuid",
			url: "/subscriptions/total?" +
				"start_date=01-2025&" +
				"end_date=12-2025&" +
				"user_id=bad-uuid",
			wantErr: true,
		},
		{
			name: "missing start date",
			url: "/subscriptions/total?" +
				"end_date=12-2025",
			wantErr: true,
		},
		{
			name: "missing end date",
			url: "/subscriptions/total?" +
				"start_date=01-2025",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				tt.url,
				nil,
			)

			filter, err := ParseQueryParams(req)

			if tt.wantErr {
				if !errors.Is(err, model.ErrValidate) {
					t.Fatalf("expected ErrValidate, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if filter == nil {
				t.Fatalf("expected filter, got nil")
			}
		})
	}
}

func TestParseQueryParams_FillFilter(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total?"+
			"start_date=01-2025&"+
			"end_date=03-2025&"+
			"user_id=550e8400-e29b-41d4-a716-446655440000&"+
			"service_name=Netflix",
		nil,
	)
	filter, err := ParseQueryParams(req)
	if err != nil {
		t.Fatal(err)
	}

	if filter.ServiceName == nil {
		t.Fatal("Service name is nil")
	}
	if *filter.ServiceName != "Netflix" {
		t.Fatalf("expected Netflix, got %s",
			*filter.ServiceName)
	}

	if filter.UserID == nil {
		t.Fatal("user id is nil")
	}
	expectedID := "550e8400-e29b-41d4-a716-446655440000"
	if filter.UserID.String() != expectedID {
		t.Fatalf("expected %s, got %s",
			expectedID,
			filter.UserID.String())
	}

	if !filter.StartDate.Equal(date(2025, time.January)) {
		t.Fatalf(
			"unexpected start date: %v",
			filter.StartDate,
		)
	}
	if !filter.EndDate.Equal(date(2025, time.March)) {
		t.Fatalf(
			"unexpected end date: %v",
			filter.EndDate,
		)
	}
}
