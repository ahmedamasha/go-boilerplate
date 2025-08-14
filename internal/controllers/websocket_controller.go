package controllers

import (
	"cusror_ai/internal/services"
	"encoding/json"
	"net/http"
)

type WebSocketController struct {
	wsService *services.WebSocketService
}

func NewWebSocketController(wsService *services.WebSocketService) *WebSocketController {
	return &WebSocketController{
		wsService: wsService,
	}
}

// HandleWebSocket handles GET /ws
func (c *WebSocketController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	c.wsService.HandleWebSocket(w, r)
}

// GetConnectionStats handles GET /api/v1/ws/stats
func (c *WebSocketController) GetConnectionStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"total_connections": c.wsService.GetConnectionCount(),
		"connected_users":   c.wsService.GetConnectedUsers(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
