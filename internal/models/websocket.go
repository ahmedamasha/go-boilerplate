package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WebSocketMessage represents a message sent through WebSocket
type WebSocketMessage struct {
	ID        uuid.UUID   `json:"id"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    *string     `json:"user_id,omitempty"`
}

// WebSocket message types
const (
	MessageTypeEventProcessed   = "event_processed"
	MessageTypeSegmentChanged   = "segment_changed"
	MessageTypeOfferGenerated   = "offer_generated"
	MessageTypeOfferUsed        = "offer_used"
	MessageTypeUserConnected    = "user_connected"
	MessageTypeUserDisconnected = "user_disconnected"
	MessageTypeError            = "error"
)

// EventProcessedData represents data for event processed notification
type EventProcessedData struct {
	Event     *Event    `json:"event"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// SegmentChangedData represents data for segment change notification
type SegmentChangedData struct {
	UserID      string    `json:"user_id"`
	OldSegments []Segment `json:"old_segments,omitempty"`
	NewSegments []Segment `json:"new_segments"`
	Timestamp   time.Time `json:"timestamp"`
}

// OfferGeneratedData represents data for offer generation notification
type OfferGeneratedData struct {
	Offer     *Offer    `json:"offer"`
	UserID    *string   `json:"user_id,omitempty"`
	SegmentID uuid.UUID `json:"segment_id"`
	Timestamp time.Time `json:"timestamp"`
}

// OfferUsedData represents data for offer usage notification
type OfferUsedData struct {
	Offer     *Offer    `json:"offer"`
	UserID    string    `json:"user_id"`
	OrderID   *string   `json:"order_id,omitempty"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// ErrorData represents error information
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewWebSocketMessage creates a new WebSocket message
func NewWebSocketMessage(messageType string, data interface{}, userID *string) *WebSocketMessage {
	return &WebSocketMessage{
		ID:        uuid.New(),
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID,
	}
}

// ToJSON converts the WebSocket message to JSON
func (m *WebSocketMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// NewEventProcessedMessage creates a message for event processed notification
func NewEventProcessedMessage(event *Event, userID string) *WebSocketMessage {
	data := &EventProcessedData{
		Event:     event,
		UserID:    userID,
		Timestamp: time.Now(),
	}
	return NewWebSocketMessage(MessageTypeEventProcessed, data, &userID)
}

// NewSegmentChangedMessage creates a message for segment change notification
func NewSegmentChangedMessage(userID string, oldSegments, newSegments []Segment) *WebSocketMessage {
	data := &SegmentChangedData{
		UserID:      userID,
		OldSegments: oldSegments,
		NewSegments: newSegments,
		Timestamp:   time.Now(),
	}
	return NewWebSocketMessage(MessageTypeSegmentChanged, data, &userID)
}

// NewOfferGeneratedMessage creates a message for offer generation notification
func NewOfferGeneratedMessage(offer *Offer, userID *string, segmentID uuid.UUID) *WebSocketMessage {
	data := &OfferGeneratedData{
		Offer:     offer,
		UserID:    userID,
		SegmentID: segmentID,
		Timestamp: time.Now(),
	}
	return NewWebSocketMessage(MessageTypeOfferGenerated, data, userID)
}

// NewOfferUsedMessage creates a message for offer usage notification
func NewOfferUsedMessage(offer *Offer, userID string, orderID *string, amount float64) *WebSocketMessage {
	data := &OfferUsedData{
		Offer:     offer,
		UserID:    userID,
		OrderID:   orderID,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	return NewWebSocketMessage(MessageTypeOfferUsed, data, &userID)
}

// NewErrorMessage creates a message for error notification
func NewErrorMessage(code, message, details string, userID *string) *WebSocketMessage {
	data := &ErrorData{
		Code:    code,
		Message: message,
		Details: details,
	}
	return NewWebSocketMessage(MessageTypeError, data, userID)
}
