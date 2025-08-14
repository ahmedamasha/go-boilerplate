package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// SegmentRules represents the JSONB rules field
type SegmentRules map[string]interface{}

// Value implements the driver.Valuer interface for SegmentRules
func (r SegmentRules) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for SegmentRules
func (r *SegmentRules) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, r)
	case string:
		return json.Unmarshal([]byte(v), r)
	default:
		return errors.New("cannot scan SegmentRules from non-string/non-bytes type")
	}
}

// Segment represents a user segment with classification rules
type Segment struct {
	ID          uuid.UUID    `json:"segment_id" db:"segment_id"`
	Name        string       `json:"segment_name" db:"segment_name"`
	Description string       `json:"description" db:"description"`
	Rules       SegmentRules `json:"rules" db:"rules"`
	IsActive    bool         `json:"is_active" db:"is_active"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// UserSegment represents the many-to-many relationship between users and segments
type UserSegment struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	SegmentID    uuid.UUID `json:"segment_id" db:"segment_id"`
	AssignedAt   time.Time `json:"assigned_at" db:"assigned_at"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	IsActive     bool      `json:"is_active" db:"is_active"`
}

// UserInterest represents interests by category or product ID for segments
type UserInterest struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	SegmentID uuid.UUID `json:"segment_id" db:"segment_id"`
	ProductID *string   `json:"product_id,omitempty" db:"product_id"`
	Category  *string   `json:"category,omitempty" db:"category"`
	Score     float64   `json:"score" db:"score"` // Interest score (0.0 to 1.0)
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Predefined segment types
const (
	SegmentHighSpender    = "High Spender"
	SegmentWindowShopper  = "Window Shopper"
	SegmentDiscountSeeker = "Discount Seeker"
	SegmentNewCustomer    = "New Customer"
	SegmentLoyal          = "Loyal Customer"
	SegmentAtRisk         = "At Risk"
)

// NewSegment creates a new segment with generated UUID
func NewSegment(name, description string, rules SegmentRules) *Segment {
	return &Segment{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Rules:       rules,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewUserSegment creates a new user segment relationship
func NewUserSegment(userID string, segmentID uuid.UUID) *UserSegment {
	return &UserSegment{
		ID:           uuid.New(),
		UserID:       userID,
		SegmentID:    segmentID,
		AssignedAt:   time.Now(),
		LastActivity: time.Now(),
		IsActive:     true,
	}
}

// NewUserInterest creates a new user interest
func NewUserInterest(userID string, segmentID uuid.UUID, productID, category *string, score float64) *UserInterest {
	return &UserInterest{
		ID:        uuid.New(),
		UserID:    userID,
		SegmentID: segmentID,
		ProductID: productID,
		Category:  category,
		Score:     score,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
