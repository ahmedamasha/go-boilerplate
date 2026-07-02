package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/services"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/refda/backend/internal/pkg/response"
)

type AdminHandler struct {
	events *services.EventService
}

func NewAdminHandler(events *services.EventService) *AdminHandler {
	return &AdminHandler{events: events}
}

func (h *AdminHandler) ReviewEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, apperrors.New("VALIDATION_ERROR", "invalid id", "معرف غير صالح", apperrors.ErrValidation))
		return
	}

	var req dto.AdminReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}

	event, err := h.events.AdminReview(c.Request.Context(), id, req.Approve)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, event)
}
