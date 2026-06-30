package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"subscription-service/internal/model"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("encode response", "error", err)
	}
}

func writeServiceError(w http.ResponseWriter, err error) {
	status := statusFromErr(err)

	if status >= http.StatusInternalServerError {
		writeError(w, status, "internal server error")
		return
	}
	writeError(w, statusFromErr(err), err.Error())
}

func statusFromErr(err error) int {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, model.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, model.ErrValidate):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func ParseID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid id %q", model.ErrValidate, idStr)
	}
	if id <= 0 {
		return 0, fmt.Errorf("%w: invalid id %q", model.ErrValidate, idStr)
	}
	return id, nil
}

func parseMonth(value string) (time.Time, error) {
	t, err := time.Parse("01-2006", value)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"%w: invalid month format %q, expected MM-YYYY",
			model.ErrValidate,
			value,
		)
	}
	return t, nil
}

func ParseQueryParams(r *http.Request) (*model.CalculateTotalCostFilter, error) {
	var filter model.CalculateTotalCostFilter
	q := r.URL.Query()
	startDateStr := q.Get("start_date")
	if startDateStr == "" {
		return nil, fmt.Errorf("%w: start_date is required", model.ErrValidate)
	}
	startDate, err := parseMonth(startDateStr)
	if err != nil {
		return nil, err
	}
	filter.StartDate = startDate

	endDateStr := q.Get("end_date")
	if endDateStr == "" {
		return nil, fmt.Errorf("%w: end_date is required", model.ErrValidate)
	}
	endDate, err := parseMonth(endDateStr)
	if err != nil {
		return nil, err
	}
	filter.EndDate = endDate

	userIDStr := q.Get("user_id")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: invalid user_id %q",
				model.ErrValidate,
				userIDStr,
			)
		}
		filter.UserID = &userID
	}
	serviceNameStr := q.Get("service_name")
	if serviceNameStr != "" {
		filter.ServiceName = &serviceNameStr
	}
	return &filter, nil
}

func ParseSubscriptionFilter(r *http.Request) (*model.SubscriptionFilter, error) {
	var filter model.SubscriptionFilter
	q := r.URL.Query()
	if userIDStr := q.Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: invalid user_id %q",
				model.ErrValidate,
				userIDStr,
			)
		}
		filter.UserID = &userID
	}
	if serviceName := q.Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}
	return &filter, nil
}
