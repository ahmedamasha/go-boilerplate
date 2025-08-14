package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// EventMetadata represents the JSONB metadata field
type EventMetadata map[string]interface{}

// Value implements the driver.Valuer interface for EventMetadata
func (m EventMetadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface for EventMetadata
func (m *EventMetadata) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	default:
		return errors.New("cannot scan EventMetadata from non-string/non-bytes type")
	}
}

// Event represents an e-commerce event
type Event struct {
	ID        uuid.UUID     `json:"event_id" db:"event_id"`
	UserID    string        `json:"user_id" db:"user_id"`
	EventType string        `json:"event_type" db:"event_type"`
	ProductID *string       `json:"product_id,omitempty" db:"product_id"`
	Category  *string       `json:"category,omitempty" db:"category"`
	Timestamp time.Time     `json:"timestamp" db:"timestamp"`
	Metadata  EventMetadata `json:"metadata" db:"metadata"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

// EventType constants
const (
	EventTypePageView      = "page_view"
	EventTypeProductView   = "product_view"
	EventTypeCartAdd       = "cart_add"
	EventTypeCartRemove    = "cart_remove"
	EventTypePurchase      = "purchase"
	EventTypeWishlistAdd   = "wishlist_add"
	EventTypeSearch        = "search"
	EventTypeCheckoutStart = "checkout_start"
)

// NewEvent creates a new event with generated UUID
func NewEvent(userID, eventType string, productID, category *string, metadata EventMetadata) *Event {
	return &Event{
		ID:        uuid.New(),
		UserID:    userID,
		EventType: eventType,
		ProductID: productID,
		Category:  category,
		Timestamp: time.Now(),
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
