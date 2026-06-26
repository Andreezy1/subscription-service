package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"subscription-service/internal/model"

	_ "subscription-service/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SubscriptionService interface {
	Create(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	GetByID(ctx context.Context, id int) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	Delete(ctx context.Context, id int) error
	CalculateTotalCost(ctx context.Context, filter model.CalculateTotalCostFilter) (int, error)
}

type Handler struct {
	service SubscriptionService
	logger  *slog.Logger
}

func NewHandler(service SubscriptionService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create godoc
//
// @Summary Создать подписку
// @Description Создает новую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
//
// @Param request body SubscriptionRequest true "subscription"
//
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
//
// @Router /subscriptions [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	log := h.log(r.Context())
	var req SubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(
			"decode create subscription request",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	subscription, err := req.ToModel()
	if err != nil {
		log.Error("convert request to model",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	createdSub, err := h.service.Create(r.Context(), subscription)
	if err != nil {
		log.Error(
			"create subscription",
			slog.String("service_name", subscription.ServiceName),
			slog.Any("error", err))
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, NewSubscriptionResponse(createdSub))
}

// Get godoc
//
// @Summary Получить подписку
// @Tags subscriptions
// @Produce json
//
// @Param id path int true "Subscription ID"
//
// @Success 200 {object} SubscriptionResponse
// @Failure 404 {object} ErrorResponse
//
// @Router /subscriptions/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	log := h.log(r.Context())

	id, err := ParseID(r)
	if err != nil {
		log.Error(
			"parse id",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Error(
			"get subscriptio",
			slog.Int("id", id),
			slog.Any("error", err))
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, NewSubscriptionResponse(sub))
}

// Delete godoc
//
// @Summary Удалить подписку
// @Tags subscriptions
//
// @Param id path int true "Subscription ID"
//
// @Success 204
// @Failure 404 {object} ErrorResponse
//
// @Router /subscriptions/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	log := h.log(r.Context())

	id, err := ParseID(r)
	if err != nil {
		log.Error(
			"parse id",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		log.Error(
			"delete subscription",
			slog.Int("id", id),
			slog.Any("error", err))
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Update godoc
//
// @Summary Обновить подписку
// @Tags subscriptions
// @Accept json
// @Produce json
//
// @Param id path int true "Subscription ID"
// @Param request body SubscriptionRequest true "subscription"
//
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
//
// @Router /subscriptions/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	log := h.log(r.Context())

	id, err := ParseID(r)
	if err != nil {
		log.Error(
			"parse id",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(
			"decode update subscription request",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sub, err := req.ToModel()
	if err != nil {
		log.Error("convert request to model",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sub.ID = id
	updatedSub, err := h.service.Update(r.Context(), sub)
	if err != nil {
		log.Error("update subscription",
			slog.Int("id", id),
			slog.Any("error", err))
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, NewSubscriptionResponse(updatedSub))
}

// CalculateTotalCost godoc
//
// @Summary Рассчитать стоимость подписок
// @Description Возвращает суммарную стоимость подписок за период
// @Tags subscriptions
// @Produce json
//
// @Param start_date query string true "MM-YYYY"
// @Param end_date query string true "MM-YYYY"
// @Param user_id query string false "UUID"
// @Param service_name query string false "Service name"
//
// @Success 200 {object} CalculateTotalCostResponse
// @Failure 400 {object} ErrorResponse
//
// @Router /subscriptions/total [get]
func (h *Handler) CalculateTotalCost(w http.ResponseWriter, r *http.Request) {
	log := h.log(r.Context())

	filter, err := ParseQueryParams(r)
	if err != nil {
		log.Error("parse query params",
			slog.Any("error", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	total, err := h.service.CalculateTotalCost(r.Context(), *filter)
	if err != nil {
		log.Error("calculate total cost",
			slog.Any("error", err))
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, NewCalculateTotalCostResponse(total))
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/subscriptions", func(r chi.Router) {
		r.Get("/total", h.CalculateTotalCost)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)

		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *Handler) log(ctx context.Context) *slog.Logger {
	reqID := GetRequestID(ctx)
	if reqID == "" {
		return h.logger
	}
	return h.logger.With(
		slog.String("request_id", reqID),
	)
}
