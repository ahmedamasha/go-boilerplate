package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/services"
	"github.com/refda/backend/internal/pkg/response"
	"github.com/refda/backend/internal/pkg/validator"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Login sends OTP to existing or new user phone
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.auth.SendOTP(c.Request.Context(), req.Phone); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "OTP sent", "message_ar": "تم إرسال رمز التحقق"})
}

// Register sends OTP for new user (requires full_name on verify)
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" validate:"required"`
		FullName string `json:"full_name" validate:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.auth.Register(c.Request.Context(), req.Phone, req.FullName); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "OTP sent", "message_ar": "تم إرسال رمز التحقق"})
}

// Verify validates OTP and returns JWT
func (h *AuthHandler) Verify(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}
	resp, err := h.auth.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, resp)
}
