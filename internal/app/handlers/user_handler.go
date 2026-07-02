package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/services"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/refda/backend/internal/pkg/middleware"
	"github.com/refda/backend/internal/pkg/response"
	"github.com/refda/backend/internal/pkg/validator"
)

type UserHandler struct {
	users *services.UserService
}

func NewUserHandler(users *services.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, apperrors.ErrUnauthorized)
		return
	}
	profile, err := h.users.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, apperrors.ErrUnauthorized)
		return
	}

	var req dto.UpdateProfileRequest
	if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
		req.FullName = c.PostForm("full_name")
		req.Location = c.PostForm("location")
	} else if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.Fail(c, err)
		return
	}

	file, _ := c.FormFile("profile_picture")

	profile, err := h.users.UpdateProfile(c.Request.Context(), userID, req, file)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profile)
}
