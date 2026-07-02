package handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/services"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/refda/backend/internal/pkg/middleware"
	"github.com/refda/backend/internal/pkg/response"
	"github.com/refda/backend/internal/pkg/validator"
)

type EventHandler struct {
	events *services.EventService
}

func NewEventHandler(events *services.EventService) *EventHandler {
	return &EventHandler{events: events}
}

func (h *EventHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, apperrors.ErrUnauthorized)
		return
	}

	var req dto.CreateEventRequest
	payload := c.PostForm("payload")
	if payload != "" {
		if err := json.Unmarshal([]byte(payload), &req); err != nil {
			response.Fail(c, err)
			return
		}
	} else if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}

	hero, _ := c.FormFile("hero_image")
	event, err := h.events.Create(c.Request.Context(), userID, req, hero)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, event)
}

func (h *EventHandler) GetList(c *gin.Context) {
	var userID *uuid.UUID
	if id, ok := middleware.GetUserID(c); ok && c.Query("mine") == "true" {
		userID = &id
	}
	events, err := h.events.List(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, events)
}

func (h *EventHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, apperrors.New("VALIDATION_ERROR", "invalid id", "معرف غير صالح", apperrors.ErrValidation))
		return
	}
	hideAmounts := c.Query("hide_amounts") == "true"
	event, err := h.events.GetByID(c.Request.Context(), id, hideAmounts)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, event)
}

func (h *EventHandler) Withdraw(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, apperrors.ErrUnauthorized)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, apperrors.New("VALIDATION_ERROR", "invalid id", "معرف غير صالح", apperrors.ErrValidation))
		return
	}

	var req dto.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}

	result, err := h.events.RequestWithdrawal(c.Request.Context(), userID, id, req.Amount)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, result)
}
