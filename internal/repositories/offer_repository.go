package repositories

import (
	"cusror_ai/internal/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type OfferRepository struct {
	db *sql.DB
}

func NewOfferRepository(db *sql.DB) *OfferRepository {
	return &OfferRepository{
		db: db,
	}
}

// CreateOffer creates a new offer in the database
func (r *OfferRepository) CreateOffer(offer *models.Offer) (*models.Offer, error) {
	query := `
		INSERT INTO offers (
			offer_id, segment_id, user_id, title, description, 
			discount_percent, discount_amount, free_shipping, min_order_value, max_discount,
			coupon_code, valid_from, valid_until, is_active, usage_limit, usage_count,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING offer_id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		offer.ID,
		offer.SegmentID,
		offer.UserID,
		offer.Title,
		offer.Description,
		offer.DiscountPercent,
		offer.DiscountAmount,
		offer.FreeShipping,
		offer.MinOrderValue,
		offer.MaxDiscount,
		offer.CouponCode,
		offer.ValidFrom,
		offer.ValidUntil,
		offer.IsActive,
		offer.UsageLimit,
		offer.UsageCount,
		offer.CreatedAt,
		offer.UpdatedAt,
	).Scan(&offer.ID, &offer.CreatedAt, &offer.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create offer: %w", err)
	}

	return offer, nil
}

// GetOfferByID retrieves an offer by its ID
func (r *OfferRepository) GetOfferByID(id uuid.UUID) (*models.Offer, error) {
	query := `
		SELECT 
			offer_id, segment_id, user_id, title, description,
			discount_percent, discount_amount, free_shipping, min_order_value, max_discount,
			coupon_code, valid_from, valid_until, is_active, usage_limit, usage_count,
			created_at, updated_at
		FROM offers
		WHERE offer_id = $1`

	var offer models.Offer
	err := r.db.QueryRow(query, id).Scan(
		&offer.ID,
		&offer.SegmentID,
		&offer.UserID,
		&offer.Title,
		&offer.Description,
		&offer.DiscountPercent,
		&offer.DiscountAmount,
		&offer.FreeShipping,
		&offer.MinOrderValue,
		&offer.MaxDiscount,
		&offer.CouponCode,
		&offer.ValidFrom,
		&offer.ValidUntil,
		&offer.IsActive,
		&offer.UsageLimit,
		&offer.UsageCount,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get offer: %w", err)
	}

	// Load applicable categories and products
	offer.ApplicableCategories, _ = r.GetOfferCategories(id)
	offer.ApplicableProducts, _ = r.GetOfferProducts(id)

	return &offer, nil
}

// GetOfferByCouponCode retrieves an offer by its coupon code
func (r *OfferRepository) GetOfferByCouponCode(couponCode string) (*models.Offer, error) {
	query := `
		SELECT 
			offer_id, segment_id, user_id, title, description,
			discount_percent, discount_amount, free_shipping, min_order_value, max_discount,
			coupon_code, valid_from, valid_until, is_active, usage_limit, usage_count,
			created_at, updated_at
		FROM offers
		WHERE coupon_code = $1`

	var offer models.Offer
	err := r.db.QueryRow(query, couponCode).Scan(
		&offer.ID,
		&offer.SegmentID,
		&offer.UserID,
		&offer.Title,
		&offer.Description,
		&offer.DiscountPercent,
		&offer.DiscountAmount,
		&offer.FreeShipping,
		&offer.MinOrderValue,
		&offer.MaxDiscount,
		&offer.CouponCode,
		&offer.ValidFrom,
		&offer.ValidUntil,
		&offer.IsActive,
		&offer.UsageLimit,
		&offer.UsageCount,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get offer by coupon code: %w", err)
	}

	// Load applicable categories and products
	offer.ApplicableCategories, _ = r.GetOfferCategories(offer.ID)
	offer.ApplicableProducts, _ = r.GetOfferProducts(offer.ID)

	return &offer, nil
}

// GetOffersBySegment retrieves all offers for a specific segment
func (r *OfferRepository) GetOffersBySegment(segmentID uuid.UUID) ([]models.Offer, error) {
	query := `
		SELECT 
			offer_id, segment_id, user_id, title, description,
			discount_percent, discount_amount, free_shipping, min_order_value, max_discount,
			coupon_code, valid_from, valid_until, is_active, usage_limit, usage_count,
			created_at, updated_at
		FROM offers
		WHERE segment_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, segmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query offers by segment: %w", err)
	}
	defer rows.Close()

	var offers []models.Offer
	for rows.Next() {
		var offer models.Offer
		err := rows.Scan(
			&offer.ID,
			&offer.SegmentID,
			&offer.UserID,
			&offer.Title,
			&offer.Description,
			&offer.DiscountPercent,
			&offer.DiscountAmount,
			&offer.FreeShipping,
			&offer.MinOrderValue,
			&offer.MaxDiscount,
			&offer.CouponCode,
			&offer.ValidFrom,
			&offer.ValidUntil,
			&offer.IsActive,
			&offer.UsageLimit,
			&offer.UsageCount,
			&offer.CreatedAt,
			&offer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer: %w", err)
		}

		// Load applicable categories and products for each offer
		offer.ApplicableCategories, _ = r.GetOfferCategories(offer.ID)
		offer.ApplicableProducts, _ = r.GetOfferProducts(offer.ID)

		offers = append(offers, offer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over offers: %w", err)
	}

	return offers, nil
}

// GetActiveOffers retrieves all currently active and valid offers
func (r *OfferRepository) GetActiveOffers() ([]models.Offer, error) {
	query := `
		SELECT 
			offer_id, segment_id, user_id, title, description,
			discount_percent, discount_amount, free_shipping, min_order_value, max_discount,
			coupon_code, valid_from, valid_until, is_active, usage_limit, usage_count,
			created_at, updated_at
		FROM offers
		WHERE is_active = true 
		  AND valid_from <= NOW() 
		  AND valid_until > NOW()
		  AND (usage_limit IS NULL OR usage_count < usage_limit)
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active offers: %w", err)
	}
	defer rows.Close()

	var offers []models.Offer
	for rows.Next() {
		var offer models.Offer
		err := rows.Scan(
			&offer.ID,
			&offer.SegmentID,
			&offer.UserID,
			&offer.Title,
			&offer.Description,
			&offer.DiscountPercent,
			&offer.DiscountAmount,
			&offer.FreeShipping,
			&offer.MinOrderValue,
			&offer.MaxDiscount,
			&offer.CouponCode,
			&offer.ValidFrom,
			&offer.ValidUntil,
			&offer.IsActive,
			&offer.UsageLimit,
			&offer.UsageCount,
			&offer.CreatedAt,
			&offer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer: %w", err)
		}

		// Load applicable categories and products for each offer
		offer.ApplicableCategories, _ = r.GetOfferCategories(offer.ID)
		offer.ApplicableProducts, _ = r.GetOfferProducts(offer.ID)

		offers = append(offers, offer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over offers: %w", err)
	}

	return offers, nil
}

// GetOffersByUser retrieves all offers for a specific user
func (r *OfferRepository) GetOffersByUser(userID string) ([]models.Offer, error) {
	query := `
		SELECT 
			offer_id, segment_id, user_id, title, description,
			discount_percent, discount_amount, free_shipping, min_order_value, max_discount,
			coupon_code, valid_from, valid_until, is_active, usage_limit, usage_count,
			created_at, updated_at
		FROM offers
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query offers by user: %w", err)
	}
	defer rows.Close()

	var offers []models.Offer
	for rows.Next() {
		var offer models.Offer
		err := rows.Scan(
			&offer.ID,
			&offer.SegmentID,
			&offer.UserID,
			&offer.Title,
			&offer.Description,
			&offer.DiscountPercent,
			&offer.DiscountAmount,
			&offer.FreeShipping,
			&offer.MinOrderValue,
			&offer.MaxDiscount,
			&offer.CouponCode,
			&offer.ValidFrom,
			&offer.ValidUntil,
			&offer.IsActive,
			&offer.UsageLimit,
			&offer.UsageCount,
			&offer.CreatedAt,
			&offer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer: %w", err)
		}

		// Load applicable categories and products for each offer
		offer.ApplicableCategories, _ = r.GetOfferCategories(offer.ID)
		offer.ApplicableProducts, _ = r.GetOfferProducts(offer.ID)

		offers = append(offers, offer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over offers: %w", err)
	}

	return offers, nil
}

// UpdateOffer updates an existing offer
func (r *OfferRepository) UpdateOffer(offer *models.Offer) (*models.Offer, error) {
	query := `
		UPDATE offers 
		SET title = $2, description = $3, discount_percent = $4, discount_amount = $5,
			free_shipping = $6, min_order_value = $7, max_discount = $8, coupon_code = $9,
			valid_from = $10, valid_until = $11, is_active = $12, usage_limit = $13,
			updated_at = $14
		WHERE offer_id = $1
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(
		query,
		offer.ID,
		offer.Title,
		offer.Description,
		offer.DiscountPercent,
		offer.DiscountAmount,
		offer.FreeShipping,
		offer.MinOrderValue,
		offer.MaxDiscount,
		offer.CouponCode,
		offer.ValidFrom,
		offer.ValidUntil,
		offer.IsActive,
		offer.UsageLimit,
		offer.UpdatedAt,
	).Scan(&offer.CreatedAt, &offer.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update offer: %w", err)
	}

	return offer, nil
}

// IncrementOfferUsage increments the usage count for an offer
func (r *OfferRepository) IncrementOfferUsage(offerID uuid.UUID) error {
	query := `
		UPDATE offers 
		SET usage_count = usage_count + 1, updated_at = NOW()
		WHERE offer_id = $1`

	result, err := r.db.Exec(query, offerID)
	if err != nil {
		return fmt.Errorf("failed to increment offer usage: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteOffer deletes an offer by its ID
func (r *OfferRepository) DeleteOffer(id uuid.UUID) error {
	query := `DELETE FROM offers WHERE offer_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// OfferCategory operations

// AddOfferCategory adds a category to an offer
func (r *OfferRepository) AddOfferCategory(offerID uuid.UUID, category string) error {
	query := `
		INSERT INTO offer_categories (id, offer_id, category)
		VALUES ($1, $2, $3)
		ON CONFLICT (offer_id, category) DO NOTHING`

	_, err := r.db.Exec(query, uuid.New(), offerID, category)
	if err != nil {
		return fmt.Errorf("failed to add offer category: %w", err)
	}

	return nil
}

// RemoveOfferCategory removes a category from an offer
func (r *OfferRepository) RemoveOfferCategory(offerID uuid.UUID, category string) error {
	query := `DELETE FROM offer_categories WHERE offer_id = $1 AND category = $2`
	_, err := r.db.Exec(query, offerID, category)
	if err != nil {
		return fmt.Errorf("failed to remove offer category: %w", err)
	}

	return nil
}

// GetOfferCategories retrieves all categories for an offer
func (r *OfferRepository) GetOfferCategories(offerID uuid.UUID) ([]string, error) {
	query := `SELECT category FROM offer_categories WHERE offer_id = $1`

	rows, err := r.db.Query(query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query offer categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over categories: %w", err)
	}

	return categories, nil
}

// OfferProduct operations

// AddOfferProduct adds a product to an offer
func (r *OfferRepository) AddOfferProduct(offerID uuid.UUID, productID string) error {
	query := `
		INSERT INTO offer_products (id, offer_id, product_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (offer_id, product_id) DO NOTHING`

	_, err := r.db.Exec(query, uuid.New(), offerID, productID)
	if err != nil {
		return fmt.Errorf("failed to add offer product: %w", err)
	}

	return nil
}

// RemoveOfferProduct removes a product from an offer
func (r *OfferRepository) RemoveOfferProduct(offerID uuid.UUID, productID string) error {
	query := `DELETE FROM offer_products WHERE offer_id = $1 AND product_id = $2`
	_, err := r.db.Exec(query, offerID, productID)
	if err != nil {
		return fmt.Errorf("failed to remove offer product: %w", err)
	}

	return nil
}

// GetOfferProducts retrieves all products for an offer
func (r *OfferRepository) GetOfferProducts(offerID uuid.UUID) ([]string, error) {
	query := `SELECT product_id FROM offer_products WHERE offer_id = $1`

	rows, err := r.db.Query(query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query offer products: %w", err)
	}
	defer rows.Close()

	var products []string
	for rows.Next() {
		var productID string
		err := rows.Scan(&productID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product ID: %w", err)
		}
		products = append(products, productID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over products: %w", err)
	}

	return products, nil
}

// OfferUsage operations

// CreateOfferUsage records the usage of an offer
func (r *OfferRepository) CreateOfferUsage(usage *models.OfferUsage) (*models.OfferUsage, error) {
	query := `
		INSERT INTO offer_usage (id, offer_id, user_id, order_id, used_at, amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, used_at`

	err := r.db.QueryRow(
		query,
		usage.ID,
		usage.OfferID,
		usage.UserID,
		usage.OrderID,
		usage.UsedAt,
		usage.Amount,
	).Scan(&usage.ID, &usage.UsedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create offer usage: %w", err)
	}

	return usage, nil
}

// GetOfferUsage retrieves usage records for an offer
func (r *OfferRepository) GetOfferUsage(offerID uuid.UUID) ([]models.OfferUsage, error) {
	query := `
		SELECT id, offer_id, user_id, order_id, used_at, amount
		FROM offer_usage
		WHERE offer_id = $1
		ORDER BY used_at DESC`

	rows, err := r.db.Query(query, offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query offer usage: %w", err)
	}
	defer rows.Close()

	var usages []models.OfferUsage
	for rows.Next() {
		var usage models.OfferUsage
		err := rows.Scan(
			&usage.ID,
			&usage.OfferID,
			&usage.UserID,
			&usage.OrderID,
			&usage.UsedAt,
			&usage.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer usage: %w", err)
		}
		usages = append(usages, usage)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over offer usage: %w", err)
	}

	return usages, nil
}

// GetUserOfferUsage retrieves usage records for a user
func (r *OfferRepository) GetUserOfferUsage(userID string) ([]models.OfferUsage, error) {
	query := `
		SELECT id, offer_id, user_id, order_id, used_at, amount
		FROM offer_usage
		WHERE user_id = $1
		ORDER BY used_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user offer usage: %w", err)
	}
	defer rows.Close()

	var usages []models.OfferUsage
	for rows.Next() {
		var usage models.OfferUsage
		err := rows.Scan(
			&usage.ID,
			&usage.OfferID,
			&usage.UserID,
			&usage.OrderID,
			&usage.UsedAt,
			&usage.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer usage: %w", err)
		}
		usages = append(usages, usage)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user offer usage: %w", err)
	}

	return usages, nil
}

// GetOfferUsageCount gets the total number of times an offer has been used
func (r *OfferRepository) GetOfferUsageCount(offerID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM offer_usage WHERE offer_id = $1`

	var count int
	err := r.db.QueryRow(query, offerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count offer usage: %w", err)
	}

	return count, nil
}
