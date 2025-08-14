package repositories

import (
	"cusror_ai/internal/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

// CreateEvent creates a new event in the database
func (r *EventRepository) CreateEvent(event *models.Event) (*models.Event, error) {
	query := `
		INSERT INTO events (event_id, user_id, event_type, product_id, category, timestamp, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING event_id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		event.ID,
		event.UserID,
		event.EventType,
		event.ProductID,
		event.Category,
		event.Timestamp,
		event.Metadata,
		event.CreatedAt,
		event.UpdatedAt,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

// GetEventByID retrieves an event by its ID
func (r *EventRepository) GetEventByID(id uuid.UUID) (*models.Event, error) {
	query := `
		SELECT event_id, user_id, event_type, product_id, category, timestamp, metadata, created_at, updated_at
		FROM events
		WHERE event_id = $1`

	var event models.Event
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.UserID,
		&event.EventType,
		&event.ProductID,
		&event.Category,
		&event.Timestamp,
		&event.Metadata,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

// GetEventsByUserID retrieves all events for a specific user
func (r *EventRepository) GetEventsByUserID(userID string, limit, offset int) ([]models.Event, error) {
	query := `
		SELECT event_id, user_id, event_type, product_id, category, timestamp, metadata, created_at, updated_at
		FROM events
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.EventType,
			&event.ProductID,
			&event.Category,
			&event.Timestamp,
			&event.Metadata,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

// GetEventsByType retrieves all events of a specific type
func (r *EventRepository) GetEventsByType(eventType string, limit, offset int) ([]models.Event, error) {
	query := `
		SELECT event_id, user_id, event_type, product_id, category, timestamp, metadata, created_at, updated_at
		FROM events
		WHERE event_type = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, eventType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query events by type: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.EventType,
			&event.ProductID,
			&event.Category,
			&event.Timestamp,
			&event.Metadata,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

// GetEventsByUserAndType retrieves events for a user filtered by event type
func (r *EventRepository) GetEventsByUserAndType(userID, eventType string, limit, offset int) ([]models.Event, error) {
	query := `
		SELECT event_id, user_id, event_type, product_id, category, timestamp, metadata, created_at, updated_at
		FROM events
		WHERE user_id = $1 AND event_type = $2
		ORDER BY timestamp DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, userID, eventType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query events by user and type: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.EventType,
			&event.ProductID,
			&event.Category,
			&event.Timestamp,
			&event.Metadata,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

// GetRecentEvents retrieves the most recent events across all users
func (r *EventRepository) GetRecentEvents(limit, offset int) ([]models.Event, error) {
	query := `
		SELECT event_id, user_id, event_type, product_id, category, timestamp, metadata, created_at, updated_at
		FROM events
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.EventType,
			&event.ProductID,
			&event.Category,
			&event.Timestamp,
			&event.Metadata,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

// CountEventsByUser counts the total number of events for a user
func (r *EventRepository) CountEventsByUser(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM events WHERE user_id = $1`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}

	return count, nil
}

// CountEventsByUserAndType counts events for a user by type
func (r *EventRepository) CountEventsByUserAndType(userID, eventType string) (int, error) {
	query := `SELECT COUNT(*) FROM events WHERE user_id = $1 AND event_type = $2`

	var count int
	err := r.db.QueryRow(query, userID, eventType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count events by type: %w", err)
	}

	return count, nil
}
