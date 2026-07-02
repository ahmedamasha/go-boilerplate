package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/services"
	"github.com/refda/backend/internal/domain"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/refda/backend/internal/pkg/middleware"
	"github.com/refda/backend/internal/pkg/response"
	"github.com/refda/backend/internal/pkg/validator"
)

type GiftHandler struct {
	gifts *services.GiftService
}

func NewGiftHandler(gifts *services.GiftService) *GiftHandler {
	return &GiftHandler{gifts: gifts}
}

func (h *GiftHandler) Contribute(c *gin.Context) {
	var req dto.ContributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}

	var gifterID *uuid.UUID
	if id, ok := middleware.GetUserID(c); ok {
		gifterID = &id
	}

	contribution, err := h.gifts.Contribute(c.Request.Context(), gifterID, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, gin.H{
		"contribution_id":  contribution.ID.String(),
		"reference_number": contribution.ReferenceNumber,
		"status":           contribution.Status,
		"amount":           contribution.Amount,
	})
}

func (h *GiftHandler) Pay(c *gin.Context) {
	var req dto.PayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}

	contributionID, err := uuid.Parse(req.ContributionID)
	if err != nil {
		response.Fail(c, apperrors.New("VALIDATION_ERROR", "invalid contribution_id", "معرف المساهمة غير صالح", apperrors.ErrValidation))
		return
	}

	result, err := h.gifts.CompletePayment(c.Request.Context(), contributionID, domain.PaymentMethod(req.PaymentMethod))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, result)
}
