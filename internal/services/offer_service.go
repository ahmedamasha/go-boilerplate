package services

import (
	"cusror_ai/internal/models"
	"cusror_ai/internal/repositories"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OfferService struct {
	offerRepo   *repositories.OfferRepository
	segmentRepo *repositories.SegmentRepository
	eventRepo   *repositories.EventRepository
	wsService   *WebSocketService
}

func NewOfferService(
	offerRepo *repositories.OfferRepository,
	segmentRepo *repositories.SegmentRepository,
	eventRepo *repositories.EventRepository,
) *OfferService {
	return &OfferService{
		offerRepo:   offerRepo,
		segmentRepo: segmentRepo,
		eventRepo:   eventRepo,
	}
}

// SetWebSocketService sets the WebSocket service to avoid circular dependencies
func (s *OfferService) SetWebSocketService(wsService *WebSocketService) {
	s.wsService = wsService
}

// CreateOffer creates a new offer
func (s *OfferService) CreateOffer(offer *models.Offer) (*models.Offer, error) {
	if offer.Title == "" {
		return nil, errors.New("offer title is required")
	}

	if offer.SegmentID == uuid.Nil {
		return nil, errors.New("segment ID is required")
	}

	// Validate segment exists
	_, err := s.segmentRepo.GetSegmentByID(offer.SegmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("segment not found")
		}
		return nil, err
	}

	// Validate offer constraints
	if err := s.validateOffer(offer); err != nil {
		return nil, err
	}

	createdOffer, err := s.offerRepo.CreateOffer(offer)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer: %w", err)
	}

	// Add categories and products if specified
	for _, category := range offer.ApplicableCategories {
		s.offerRepo.AddOfferCategory(createdOffer.ID, category)
	}

	for _, productID := range offer.ApplicableProducts {
		s.offerRepo.AddOfferProduct(createdOffer.ID, productID)
	}

	// Send WebSocket notification
	if s.wsService != nil {
		s.wsService.NotifyOfferGenerated(createdOffer, offer.UserID, offer.SegmentID)
	}

	return createdOffer, nil
}

// GetOfferByID retrieves an offer by its ID
func (s *OfferService) GetOfferByID(id uuid.UUID) (*models.Offer, error) {
	offer, err := s.offerRepo.GetOfferByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("offer not found")
		}
		return nil, err
	}
	return offer, nil
}

// GetOfferByCouponCode retrieves an offer by its coupon code
func (s *OfferService) GetOfferByCouponCode(couponCode string) (*models.Offer, error) {
	offer, err := s.offerRepo.GetOfferByCouponCode(couponCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("offer not found")
		}
		return nil, err
	}
	return offer, nil
}

// GetActiveOffers retrieves all currently active offers
func (s *OfferService) GetActiveOffers() ([]models.Offer, error) {
	return s.offerRepo.GetActiveOffers()
}

// GetOffersByUser retrieves all offers for a specific user
func (s *OfferService) GetOffersByUser(userID string) ([]models.Offer, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.offerRepo.GetOffersByUser(userID)
}

// GetOffersBySegment retrieves all offers for a specific segment
func (s *OfferService) GetOffersBySegment(segmentID uuid.UUID) ([]models.Offer, error) {
	return s.offerRepo.GetOffersBySegment(segmentID)
}

// UseOffer records the usage of an offer
func (s *OfferService) UseOffer(offerID uuid.UUID, userID string, orderID *string, amount float64) (*models.OfferUsage, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if amount < 0 {
		return nil, errors.New("amount must be non-negative")
	}

	// Get the offer
	offer, err := s.offerRepo.GetOfferByID(offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer: %w", err)
	}

	// Validate offer can be used
	if !offer.IsValid() {
		return nil, errors.New("offer is not valid")
	}

	if !offer.CanBeUsedBy(userID) {
		return nil, errors.New("offer cannot be used by this user")
	}

	// Create usage record
	usage := &models.OfferUsage{
		ID:      uuid.New(),
		OfferID: offerID,
		UserID:  userID,
		OrderID: orderID,
		UsedAt:  time.Now(),
		Amount:  amount,
	}

	createdUsage, err := s.offerRepo.CreateOfferUsage(usage)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer usage: %w", err)
	}

	// Increment offer usage count
	err = s.offerRepo.IncrementOfferUsage(offerID)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to increment offer usage count: %v\n", err)
	}

	// Send WebSocket notification
	if s.wsService != nil {
		s.wsService.NotifyOfferUsed(offer, userID, orderID, amount)
	}

	return createdUsage, nil
}

// GenerateOffersForUserSegment generates personalized offers for a user based on their segment
func (s *OfferService) GenerateOffersForUserSegment(userID string, segmentID uuid.UUID) ([]*models.Offer, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get segment information
	segment, err := s.segmentRepo.GetSegmentByID(segmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get segment: %w", err)
	}

	// Get user interests to personalize offers
	interests, err := s.segmentRepo.GetUserInterests(userID)
	if err != nil {
		interests = []models.UserInterest{} // Continue with empty interests
	}

	// Generate offers based on segment
	var offers []*models.Offer

	switch segment.Name {
	case models.SegmentHighSpender:
		offers = s.generateHighSpenderOffers(userID, segmentID, interests)
	case models.SegmentWindowShopper:
		offers = s.generateWindowShopperOffers(userID, segmentID, interests)
	case models.SegmentDiscountSeeker:
		offers = s.generateDiscountSeekerOffers(userID, segmentID, interests)
	default:
		offers = s.generateGenericOffers(userID, segmentID, interests)
	}

	// Create offers in database
	var createdOffers []*models.Offer
	for _, offer := range offers {
		createdOffer, err := s.CreateOffer(offer)
		if err != nil {
			// Log error but continue with other offers
			fmt.Printf("Warning: failed to create offer: %v\n", err)
			continue
		}
		createdOffers = append(createdOffers, createdOffer)
	}

	return createdOffers, nil
}

// generateHighSpenderOffers generates offers for high-spending users
func (s *OfferService) generateHighSpenderOffers(userID string, segmentID uuid.UUID, interests []models.UserInterest) []*models.Offer {
	var offers []*models.Offer

	// VIP exclusive discount
	vipOffer := models.NewPercentageOffer(
		segmentID,
		&userID,
		"VIP Exclusive: 15% Off Premium Items",
		"Exclusive discount for our valued premium customers",
		15.0,
		nil,                  // No max discount for VIP
		&[]float64{100.0}[0], // Min order $100
	)
	vipOffer.GenerateCouponCode()
	offers = append(offers, vipOffer)

	// Free shipping on any order
	freeShippingOffer := models.NewFreeShippingOffer(
		segmentID,
		&userID,
		"Free Premium Shipping",
		"Complimentary premium shipping on all orders",
		nil, // No minimum order value
	)
	freeShippingOffer.GenerateCouponCode()
	offers = append(offers, freeShippingOffer)

	return offers
}

// generateWindowShopperOffers generates offers for window shoppers
func (s *OfferService) generateWindowShopperOffers(userID string, segmentID uuid.UUID, interests []models.UserInterest) []*models.Offer {
	var offers []*models.Offer

	// Limited time percentage discount to encourage purchase
	urgencyOffer := models.NewPercentageOffer(
		segmentID,
		&userID,
		"Limited Time: 20% Off - Today Only!",
		"Don't miss out! Special discount expires in 24 hours",
		20.0,
		&[]float64{50.0}[0], // Max discount $50
		&[]float64{25.0}[0], // Min order $25
	)
	urgencyOffer.ValidUntil = time.Now().Add(24 * time.Hour) // 24 hour expiry
	urgencyOffer.GenerateCouponCode()
	offers = append(offers, urgencyOffer)

	// Free shipping to reduce barriers
	freeShippingOffer := models.NewFreeShippingOffer(
		segmentID,
		&userID,
		"Free Shipping - No Minimum!",
		"Free shipping on any order, no minimum required",
		nil,
	)
	freeShippingOffer.GenerateCouponCode()
	offers = append(offers, freeShippingOffer)

	return offers
}

// generateDiscountSeekerOffers generates offers for discount seekers
func (s *OfferService) generateDiscountSeekerOffers(userID string, segmentID uuid.UUID, interests []models.UserInterest) []*models.Offer {
	var offers []*models.Offer

	// High percentage discount
	bigDiscountOffer := models.NewPercentageOffer(
		segmentID,
		&userID,
		"MEGA DEAL: 30% Off Everything!",
		"Huge savings on all items - limited time offer",
		30.0,
		&[]float64{100.0}[0], // Max discount $100
		&[]float64{50.0}[0],  // Min order $50
	)
	bigDiscountOffer.GenerateCouponCode()
	offers = append(offers, bigDiscountOffer)

	// Bundle offer with higher perceived value
	bundleOffer := models.NewFixedAmountOffer(
		segmentID,
		&userID,
		"Buy More, Save More: $25 Off Orders Over $100",
		"Get $25 off when you spend $100 or more",
		&[]float64{25.0}[0],  // $25 off
		&[]float64{100.0}[0], // Min order $100
	)
	bundleOffer.GenerateCouponCode()
	offers = append(offers, bundleOffer)

	return offers
}

// generateGenericOffers generates general offers for other segments
func (s *OfferService) generateGenericOffers(userID string, segmentID uuid.UUID, interests []models.UserInterest) []*models.Offer {
	var offers []*models.Offer

	// Standard percentage discount
	standardOffer := models.NewPercentageOffer(
		segmentID,
		&userID,
		"Welcome Back: 10% Off Your Next Purchase",
		"Enjoy 10% off your next order as a valued customer",
		10.0,
		&[]float64{25.0}[0], // Max discount $25
		&[]float64{30.0}[0], // Min order $30
	)
	standardOffer.GenerateCouponCode()
	offers = append(offers, standardOffer)

	return offers
}

// ValidateOffer validates offer constraints
func (s *OfferService) validateOffer(offer *models.Offer) error {
	// Check discount constraints
	if offer.DiscountPercent != nil {
		if *offer.DiscountPercent < 0 || *offer.DiscountPercent > 100 {
			return errors.New("discount percentage must be between 0 and 100")
		}
	}

	if offer.DiscountAmount != nil {
		if *offer.DiscountAmount < 0 {
			return errors.New("discount amount must be non-negative")
		}
	}

	// Check that at least one discount type is specified
	if offer.DiscountPercent == nil && offer.DiscountAmount == nil && !offer.FreeShipping {
		return errors.New("offer must have at least one discount type")
	}

	// Check date constraints
	if offer.ValidUntil.Before(offer.ValidFrom) {
		return errors.New("valid until date must be after valid from date")
	}

	// Check usage limit
	if offer.UsageLimit != nil && *offer.UsageLimit < 0 {
		return errors.New("usage limit must be non-negative")
	}

	return nil
}

// GetOfferUsage retrieves usage records for an offer
func (s *OfferService) GetOfferUsage(offerID uuid.UUID) ([]models.OfferUsage, error) {
	return s.offerRepo.GetOfferUsage(offerID)
}

// GetUserOfferUsage retrieves usage records for a user
func (s *OfferService) GetUserOfferUsage(userID string) ([]models.OfferUsage, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.offerRepo.GetUserOfferUsage(userID)
}

// UpdateOffer updates an existing offer
func (s *OfferService) UpdateOffer(offer *models.Offer) (*models.Offer, error) {
	if offer.Title == "" {
		return nil, errors.New("offer title is required")
	}

	// Validate offer constraints
	if err := s.validateOffer(offer); err != nil {
		return nil, err
	}

	offer.UpdatedAt = time.Now()

	updatedOffer, err := s.offerRepo.UpdateOffer(offer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("offer not found")
		}
		return nil, err
	}

	return updatedOffer, nil
}

// DeleteOffer deletes an offer
func (s *OfferService) DeleteOffer(id uuid.UUID) error {
	err := s.offerRepo.DeleteOffer(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("offer not found")
		}
		return err
	}
	return nil
}

// GetOfferStats retrieves statistics about offers
func (s *OfferService) GetOfferStats() (*OfferStats, error) {
	activeOffers, err := s.offerRepo.GetActiveOffers()
	if err != nil {
		return nil, err
	}

	stats := &OfferStats{
		TotalActiveOffers: len(activeOffers),
		OffersByType:      make(map[string]int),
	}

	for _, offer := range activeOffers {
		if offer.DiscountPercent != nil {
			stats.OffersByType["percentage"]++
		} else if offer.DiscountAmount != nil {
			stats.OffersByType["fixed_amount"]++
		} else if offer.FreeShipping {
			stats.OffersByType["free_shipping"]++
		}
	}

	return stats, nil
}

// OfferStats represents statistics about offers
type OfferStats struct {
	TotalActiveOffers int            `json:"total_active_offers"`
	OffersByType      map[string]int `json:"offers_by_type"`
}
