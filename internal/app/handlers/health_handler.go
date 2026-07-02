package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/refda/backend/internal/pkg/response"
)

func Health(c *gin.Context) {
	response.OK(c, gin.H{
		"status":  "ok",
		"service": "refda-api",
		"version": "1.0.0",
	})
}
