package controllers

import (
	"cusror_ai/internal/models"
	"cusror_ai/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EventController struct {
	eventService *services.EventService
}

func NewEventController(eventService *services.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

// CreateEvent handles POST /api/v1/events
func (c *EventController) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" || req.EventType == "" {
		http.Error(w, "user_id and event_type are required", http.StatusBadRequest)
		return
	}

	// Validate event type
	if err := c.eventService.ValidateEventType(req.EventType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create event
	event, err := c.eventService.CreateEvent(
		req.UserID,
		req.EventType,
		req.ProductID,
		req.Category,
		req.Metadata,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// GetEvent handles GET /api/v1/events/{id}
func (c *EventController) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := c.eventService.GetEventByID(id)
	if err != nil {
		if err.Error() == "event not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// GetUserEvents handles GET /api/v1/events/user/{user_id}
func (c *EventController) GetUserEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	events, err := c.eventService.GetEventsByUser(userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
		"limit":  limit,
		"offset": offset,
		"count":  len(events),
	})
}

// GetEventsByType handles GET /api/v1/events/type/{type}
func (c *EventController) GetEventsByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventType := vars["type"]

	// Validate event type
	if err := c.eventService.ValidateEventType(eventType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	events, err := c.eventService.GetEventsByType(eventType, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events":     events,
		"event_type": eventType,
		"limit":      limit,
		"offset":     offset,
		"count":      len(events),
	})
}

// GetRecentEvents handles GET /api/v1/events/recent
func (c *EventController) GetRecentEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	events, err := c.eventService.GetRecentEvents(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
		"limit":  limit,
		"offset": offset,
		"count":  len(events),
	})
}

// GetUserEventStats handles GET /api/v1/events/user/{user_id}/stats
func (c *EventController) GetUserEventStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	stats, err := c.eventService.GetUserEventStats(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// SimulateShopifyEvent handles POST /api/v1/events/simulate/shopify
func (c *EventController) SimulateShopifyEvent(w http.ResponseWriter, r *http.Request) {
	var req ShopifyWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert Shopify webhook to internal event
	event, err := c.convertShopifyToEvent(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the event
	createdEvent, err := c.eventService.CreateEvent(
		event.UserID,
		event.EventType,
		event.ProductID,
		event.Category,
		event.Metadata,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Shopify event simulated successfully",
		"event":   createdEvent,
	})
}

// convertShopifyToEvent converts a simulated Shopify webhook to internal event
func (c *EventController) convertShopifyToEvent(req ShopifyWebhookRequest) (*EventData, error) {
	event := &EventData{
		UserID:   req.CustomerID,
		Metadata: make(models.EventMetadata),
	}

	// Map Shopify events to internal event types
	switch req.EventType {
	case "order_created":
		event.EventType = models.EventTypePurchase
		if req.OrderData != nil {
			event.Metadata["amount"] = req.OrderData.TotalPrice
			event.Metadata["order_id"] = req.OrderData.OrderID
			event.Metadata["currency"] = req.OrderData.Currency
		}
	case "product_viewed":
		event.EventType = models.EventTypeProductView
		if req.ProductData != nil {
			event.ProductID = &req.ProductData.ProductID
			event.Category = &req.ProductData.Category
			event.Metadata["product_name"] = req.ProductData.Title
			event.Metadata["price"] = req.ProductData.Price
		}
	case "cart_updated":
		event.EventType = models.EventTypeCartAdd
		if req.ProductData != nil {
			event.ProductID = &req.ProductData.ProductID
			event.Category = &req.ProductData.Category
			event.Metadata["quantity"] = req.ProductData.Quantity
		}
	case "page_viewed":
		event.EventType = models.EventTypePageView
		event.Metadata["page_url"] = req.PageURL
	default:
		return nil, fmt.Errorf("unsupported Shopify event type: %s", req.EventType)
	}

	return event, nil
}

// Request/Response structures

type CreateEventRequest struct {
	UserID    string               `json:"user_id"`
	EventType string               `json:"event_type"`
	ProductID *string              `json:"product_id,omitempty"`
	Category  *string              `json:"category,omitempty"`
	Metadata  models.EventMetadata `json:"metadata,omitempty"`
}

type EventData struct {
	UserID    string               `json:"user_id"`
	EventType string               `json:"event_type"`
	ProductID *string              `json:"product_id,omitempty"`
	Category  *string              `json:"category,omitempty"`
	Metadata  models.EventMetadata `json:"metadata,omitempty"`
}

// Shopify simulation structures
type ShopifyWebhookRequest struct {
	EventType   string       `json:"event_type"`
	CustomerID  string       `json:"customer_id"`
	PageURL     string       `json:"page_url,omitempty"`
	OrderData   *OrderData   `json:"order_data,omitempty"`
	ProductData *ProductData `json:"product_data,omitempty"`
}

type OrderData struct {
	OrderID    string  `json:"order_id"`
	TotalPrice float64 `json:"total_price"`
	Currency   string  `json:"currency"`
}

type ProductData struct {
	ProductID string  `json:"product_id"`
	Title     string  `json:"title"`
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity,omitempty"`
}
