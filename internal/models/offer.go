package models

import (
	"time"

	"github.com/google/uuid"
)

// Offer represents a personalized offer or coupon
type Offer struct {
	ID                   uuid.UUID `json:"offer_id" db:"offer_id"`
	SegmentID            uuid.UUID `json:"segment_id" db:"segment_id"`
	UserID               *string   `json:"user_id,omitempty" db:"user_id"` // Specific user or null for segment-wide
	Title                string    `json:"title" db:"title"`
	Description          string    `json:"description" db:"description"`
	DiscountPercent      *float64  `json:"discount_percent,omitempty" db:"discount_percent"`
	DiscountAmount       *float64  `json:"discount_amount,omitempty" db:"discount_amount"`
	FreeShipping         bool      `json:"free_shipping" db:"free_shipping"`
	MinOrderValue        *float64  `json:"min_order_value,omitempty" db:"min_order_value"`
	MaxDiscount          *float64  `json:"max_discount,omitempty" db:"max_discount"`
	CouponCode           *string   `json:"coupon_code,omitempty" db:"coupon_code"`
	ValidFrom            time.Time `json:"valid_from" db:"valid_from"`
	ValidUntil           time.Time `json:"valid_until" db:"valid_until"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	UsageLimit           *int      `json:"usage_limit,omitempty" db:"usage_limit"`
	UsageCount           int       `json:"usage_count" db:"usage_count"`
	ApplicableCategories []string  `json:"applicable_categories,omitempty" db:"-"` // Will be stored in separate table
	ApplicableProducts   []string  `json:"applicable_products,omitempty" db:"-"`   // Will be stored in separate table
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// OfferCategory represents categories where an offer can be applied
type OfferCategory struct {
	ID       uuid.UUID `json:"id" db:"id"`
	OfferID  uuid.UUID `json:"offer_id" db:"offer_id"`
	Category string    `json:"category" db:"category"`
}

// OfferProduct represents specific products where an offer can be applied
type OfferProduct struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OfferID   uuid.UUID `json:"offer_id" db:"offer_id"`
	ProductID string    `json:"product_id" db:"product_id"`
}

// OfferUsage tracks when and how offers are used
type OfferUsage struct {
	ID      uuid.UUID `json:"id" db:"id"`
	OfferID uuid.UUID `json:"offer_id" db:"offer_id"`
	UserID  string    `json:"user_id" db:"user_id"`
	OrderID *string   `json:"order_id,omitempty" db:"order_id"`
	UsedAt  time.Time `json:"used_at" db:"used_at"`
	Amount  float64   `json:"amount" db:"amount"` // Discount amount applied
}

// Offer type constants
const (
	OfferTypePercentage   = "percentage"
	OfferTypeFixedAmount  = "fixed_amount"
	OfferTypeFreeShipping = "free_shipping"
	OfferTypeBuyOneGetOne = "bogo"
)

// NewOffer creates a new offer with generated UUID
func NewOffer(segmentID uuid.UUID, userID *string, title, description string) *Offer {
	now := time.Now()
	return &Offer{
		ID:          uuid.New(),
		SegmentID:   segmentID,
		UserID:      userID,
		Title:       title,
		Description: description,
		ValidFrom:   now,
		ValidUntil:  now.Add(30 * 24 * time.Hour), // Default 30 days validity
		IsActive:    true,
		UsageCount:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewPercentageOffer creates a percentage-based offer
func NewPercentageOffer(segmentID uuid.UUID, userID *string, title, description string,
	discountPercent float64, maxDiscount *float64, minOrderValue *float64) *Offer {
	offer := NewOffer(segmentID, userID, title, description)
	offer.DiscountPercent = &discountPercent
	offer.MaxDiscount = maxDiscount
	offer.MinOrderValue = minOrderValue
	return offer
}

// NewFixedAmountOffer creates a fixed amount discount offer
func NewFixedAmountOffer(segmentID uuid.UUID, userID *string, title, description string,
	discountAmount, minOrderValue *float64) *Offer {
	offer := NewOffer(segmentID, userID, title, description)
	offer.DiscountAmount = discountAmount
	offer.MinOrderValue = minOrderValue
	return offer
}

// NewFreeShippingOffer creates a free shipping offer
func NewFreeShippingOffer(segmentID uuid.UUID, userID *string, title, description string,
	minOrderValue *float64) *Offer {
	offer := NewOffer(segmentID, userID, title, description)
	offer.FreeShipping = true
	offer.MinOrderValue = minOrderValue
	return offer
}

// IsValid checks if the offer is currently valid
func (o *Offer) IsValid() bool {
	now := time.Now()
	return o.IsActive &&
		now.After(o.ValidFrom) &&
		now.Before(o.ValidUntil) &&
		(o.UsageLimit == nil || o.UsageCount < *o.UsageLimit)
}

// CanBeUsedBy checks if the offer can be used by a specific user
func (o *Offer) CanBeUsedBy(userID string) bool {
	if !o.IsValid() {
		return false
	}
	// If offer is user-specific, check if it matches
	if o.UserID != nil {
		return *o.UserID == userID
	}
	// Segment-wide offer can be used by any user in the segment
	return true
}

// GenerateCouponCode generates a simple coupon code for the offer
func (o *Offer) GenerateCouponCode() string {
	// Simple implementation - in production, use a more sophisticated generator
	code := "OFFER" + o.ID.String()[:8]
	o.CouponCode = &code
	return code
}
