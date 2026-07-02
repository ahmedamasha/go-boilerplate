package response

import (
	"net/http"

	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	MessageAR string `json:"message_ar,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Success: true, Data: data})
}

func Fail(c *gin.Context, err error) {
	appErr := apperrors.Map(err)
	status := http.StatusInternalServerError

	switch appErr.Code {
	case "NOT_FOUND":
		status = http.StatusNotFound
	case "UNAUTHORIZED", "INVALID_OTP", "OTP_EXPIRED":
		status = http.StatusUnauthorized
	case "FORBIDDEN":
		status = http.StatusForbidden
	case "VALIDATION_ERROR", "INVALID_PHONE", "USER_EXISTS", "GIFT_RESERVED", "INSUFFICIENT_AMOUNT", "EVENT_NOT_OPEN", "EVENT_NOT_WITHDRAWABLE", "INVALID_WITHDRAWAL":
		status = http.StatusBadRequest
	}

	c.JSON(status, Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:      appErr.Code,
			Message:   appErr.Message,
			MessageAR: appErr.MessageAR,
		},
	})
}
