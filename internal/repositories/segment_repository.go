package repositories

import (
	"cusror_ai/internal/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type SegmentRepository struct {
	db *sql.DB
}

func NewSegmentRepository(db *sql.DB) *SegmentRepository {
	return &SegmentRepository{
		db: db,
	}
}

// CreateSegment creates a new segment in the database
func (r *SegmentRepository) CreateSegment(segment *models.Segment) (*models.Segment, error) {
	query := `
		INSERT INTO segments (segment_id, segment_name, description, rules, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING segment_id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		segment.ID,
		segment.Name,
		segment.Description,
		segment.Rules,
		segment.IsActive,
		segment.CreatedAt,
		segment.UpdatedAt,
	).Scan(&segment.ID, &segment.CreatedAt, &segment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create segment: %w", err)
	}

	return segment, nil
}

// GetSegmentByID retrieves a segment by its ID
func (r *SegmentRepository) GetSegmentByID(id uuid.UUID) (*models.Segment, error) {
	query := `
		SELECT segment_id, segment_name, description, rules, is_active, created_at, updated_at
		FROM segments
		WHERE segment_id = $1`

	var segment models.Segment
	err := r.db.QueryRow(query, id).Scan(
		&segment.ID,
		&segment.Name,
		&segment.Description,
		&segment.Rules,
		&segment.IsActive,
		&segment.CreatedAt,
		&segment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get segment: %w", err)
	}

	return &segment, nil
}

// GetSegmentByName retrieves a segment by its name
func (r *SegmentRepository) GetSegmentByName(name string) (*models.Segment, error) {
	query := `
		SELECT segment_id, segment_name, description, rules, is_active, created_at, updated_at
		FROM segments
		WHERE segment_name = $1`

	var segment models.Segment
	err := r.db.QueryRow(query, name).Scan(
		&segment.ID,
		&segment.Name,
		&segment.Description,
		&segment.Rules,
		&segment.IsActive,
		&segment.CreatedAt,
		&segment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get segment by name: %w", err)
	}

	return &segment, nil
}

// GetAllSegments retrieves all segments
func (r *SegmentRepository) GetAllSegments() ([]models.Segment, error) {
	query := `
		SELECT segment_id, segment_name, description, rules, is_active, created_at, updated_at
		FROM segments
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query segments: %w", err)
	}
	defer rows.Close()

	var segments []models.Segment
	for rows.Next() {
		var segment models.Segment
		err := rows.Scan(
			&segment.ID,
			&segment.Name,
			&segment.Description,
			&segment.Rules,
			&segment.IsActive,
			&segment.CreatedAt,
			&segment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan segment: %w", err)
		}
		segments = append(segments, segment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over segments: %w", err)
	}

	return segments, nil
}

// GetActiveSegments retrieves all active segments
func (r *SegmentRepository) GetActiveSegments() ([]models.Segment, error) {
	query := `
		SELECT segment_id, segment_name, description, rules, is_active, created_at, updated_at
		FROM segments
		WHERE is_active = true
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active segments: %w", err)
	}
	defer rows.Close()

	var segments []models.Segment
	for rows.Next() {
		var segment models.Segment
		err := rows.Scan(
			&segment.ID,
			&segment.Name,
			&segment.Description,
			&segment.Rules,
			&segment.IsActive,
			&segment.CreatedAt,
			&segment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan segment: %w", err)
		}
		segments = append(segments, segment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over segments: %w", err)
	}

	return segments, nil
}

// UpdateSegment updates an existing segment
func (r *SegmentRepository) UpdateSegment(segment *models.Segment) (*models.Segment, error) {
	query := `
		UPDATE segments 
		SET segment_name = $2, description = $3, rules = $4, is_active = $5, updated_at = $6
		WHERE segment_id = $1
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(
		query,
		segment.ID,
		segment.Name,
		segment.Description,
		segment.Rules,
		segment.IsActive,
		segment.UpdatedAt,
	).Scan(&segment.CreatedAt, &segment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update segment: %w", err)
	}

	return segment, nil
}

// DeleteSegment deletes a segment by its ID
func (r *SegmentRepository) DeleteSegment(id uuid.UUID) error {
	query := `DELETE FROM segments WHERE segment_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete segment: %w", err)
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

// UserSegment operations

// AssignUserToSegment assigns a user to a segment
func (r *SegmentRepository) AssignUserToSegment(userSegment *models.UserSegment) (*models.UserSegment, error) {
	query := `
		INSERT INTO user_segments (id, user_id, segment_id, assigned_at, last_activity, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, segment_id) 
		DO UPDATE SET is_active = $6, last_activity = $5
		RETURNING id, assigned_at`

	err := r.db.QueryRow(
		query,
		userSegment.ID,
		userSegment.UserID,
		userSegment.SegmentID,
		userSegment.AssignedAt,
		userSegment.LastActivity,
		userSegment.IsActive,
	).Scan(&userSegment.ID, &userSegment.AssignedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to assign user to segment: %w", err)
	}

	return userSegment, nil
}

// RemoveUserFromSegment removes a user from a segment
func (r *SegmentRepository) RemoveUserFromSegment(userID string, segmentID uuid.UUID) error {
	query := `
		UPDATE user_segments 
		SET is_active = false 
		WHERE user_id = $1 AND segment_id = $2`

	result, err := r.db.Exec(query, userID, segmentID)
	if err != nil {
		return fmt.Errorf("failed to remove user from segment: %w", err)
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

// GetUserSegments retrieves all active segments for a user
func (r *SegmentRepository) GetUserSegments(userID string) ([]models.Segment, error) {
	query := `
		SELECT s.segment_id, s.segment_name, s.description, s.rules, s.is_active, s.created_at, s.updated_at
		FROM segments s
		INNER JOIN user_segments us ON s.segment_id = us.segment_id
		WHERE us.user_id = $1 AND us.is_active = true AND s.is_active = true
		ORDER BY us.assigned_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user segments: %w", err)
	}
	defer rows.Close()

	var segments []models.Segment
	for rows.Next() {
		var segment models.Segment
		err := rows.Scan(
			&segment.ID,
			&segment.Name,
			&segment.Description,
			&segment.Rules,
			&segment.IsActive,
			&segment.CreatedAt,
			&segment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user segment: %w", err)
		}
		segments = append(segments, segment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user segments: %w", err)
	}

	return segments, nil
}

// GetSegmentUsers retrieves all users in a segment
func (r *SegmentRepository) GetSegmentUsers(segmentID uuid.UUID) ([]string, error) {
	query := `
		SELECT user_id
		FROM user_segments
		WHERE segment_id = $1 AND is_active = true
		ORDER BY assigned_at DESC`

	rows, err := r.db.Query(query, segmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query segment users: %w", err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over segment users: %w", err)
	}

	return userIDs, nil
}

// IsUserInSegment checks if a user is in a specific segment
func (r *SegmentRepository) IsUserInSegment(userID string, segmentID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_segments 
			WHERE user_id = $1 AND segment_id = $2 AND is_active = true
		)`

	var exists bool
	err := r.db.QueryRow(query, userID, segmentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user segment: %w", err)
	}

	return exists, nil
}

// UserInterest operations

// CreateUserInterest creates a new user interest
func (r *SegmentRepository) CreateUserInterest(interest *models.UserInterest) (*models.UserInterest, error) {
	query := `
		INSERT INTO user_interests (id, user_id, segment_id, product_id, category, score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		interest.ID,
		interest.UserID,
		interest.SegmentID,
		interest.ProductID,
		interest.Category,
		interest.Score,
		interest.CreatedAt,
		interest.UpdatedAt,
	).Scan(&interest.ID, &interest.CreatedAt, &interest.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user interest: %w", err)
	}

	return interest, nil
}

// GetUserInterests retrieves all interests for a user
func (r *SegmentRepository) GetUserInterests(userID string) ([]models.UserInterest, error) {
	query := `
		SELECT id, user_id, segment_id, product_id, category, score, created_at, updated_at
		FROM user_interests
		WHERE user_id = $1
		ORDER BY score DESC, created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user interests: %w", err)
	}
	defer rows.Close()

	var interests []models.UserInterest
	for rows.Next() {
		var interest models.UserInterest
		err := rows.Scan(
			&interest.ID,
			&interest.UserID,
			&interest.SegmentID,
			&interest.ProductID,
			&interest.Category,
			&interest.Score,
			&interest.CreatedAt,
			&interest.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user interest: %w", err)
		}
		interests = append(interests, interest)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user interests: %w", err)
	}

	return interests, nil
}

// UpdateUserInterestScore updates the score for a user interest
func (r *SegmentRepository) UpdateUserInterestScore(userID string, segmentID uuid.UUID, productID, category *string, newScore float64) error {
	query := `
		UPDATE user_interests 
		SET score = $4, updated_at = NOW()
		WHERE user_id = $1 AND segment_id = $2 
		  AND (product_id = $3 OR (product_id IS NULL AND $3 IS NULL))
		  AND (category = $5 OR (category IS NULL AND $5 IS NULL))`

	result, err := r.db.Exec(query, userID, segmentID, productID, newScore, category)
	if err != nil {
		return fmt.Errorf("failed to update user interest score: %w", err)
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
