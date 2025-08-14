package services

import (
	"cusror_ai/internal/models"
	"cusror_ai/internal/repositories"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type EventService struct {
	eventRepo      *repositories.EventRepository
	segmentRepo    *repositories.SegmentRepository
	offerRepo      *repositories.OfferRepository
	segmentService *SegmentService
	offerService   *OfferService
	wsService      *WebSocketService
}

func NewEventService(
	eventRepo *repositories.EventRepository,
	segmentRepo *repositories.SegmentRepository,
	offerRepo *repositories.OfferRepository,
) *EventService {
	return &EventService{
		eventRepo:   eventRepo,
		segmentRepo: segmentRepo,
		offerRepo:   offerRepo,
	}
}

// SetServices sets the dependent services to avoid circular dependencies
func (s *EventService) SetServices(segmentService *SegmentService, offerService *OfferService, wsService *WebSocketService) {
	s.segmentService = segmentService
	s.offerService = offerService
	s.wsService = wsService
}

// CreateEvent creates a new event and triggers segmentation analysis
func (s *EventService) CreateEvent(userID, eventType string, productID, category *string, metadata models.EventMetadata) (*models.Event, error) {
	if userID == "" || eventType == "" {
		return nil, errors.New("user ID and event type are required")
	}

	// Create the event
	event := models.NewEvent(userID, eventType, productID, category, metadata)

	createdEvent, err := s.eventRepo.CreateEvent(event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Process the event asynchronously
	go s.processEventAsync(createdEvent)

	return createdEvent, nil
}

// processEventAsync handles event processing in background
func (s *EventService) processEventAsync(event *models.Event) {
	// Get current user segments
	oldSegments, _ := s.segmentRepo.GetUserSegments(event.UserID)

	// Trigger segmentation analysis
	if s.segmentService != nil {
		newSegments, err := s.segmentService.AnalyzeUserSegments(event.UserID)
		if err == nil {
			// Check if segments changed
			if s.segmentsChanged(oldSegments, newSegments) {
				// Update user interests based on the event
				s.updateUserInterests(event, newSegments)

				// Generate offers for new segments
				if s.offerService != nil {
					for _, segment := range newSegments {
						s.offerService.GenerateOffersForUserSegment(event.UserID, segment.ID)
					}
				}

				// Send WebSocket notification about segment change
				if s.wsService != nil {
					s.wsService.NotifySegmentChanged(event.UserID, oldSegments, newSegments)
				}
			}
		}
	}

	// Send WebSocket notification about event processing
	if s.wsService != nil {
		s.wsService.NotifyEventProcessed(event, event.UserID)
	}
}

// GetEventByID retrieves an event by its ID
func (s *EventService) GetEventByID(id uuid.UUID) (*models.Event, error) {
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return event, nil
}

// GetEventsByUser retrieves events for a specific user with pagination
func (s *EventService) GetEventsByUser(userID string, limit, offset int) ([]models.Event, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if limit <= 0 {
		limit = 20 // Default limit
	}

	return s.eventRepo.GetEventsByUserID(userID, limit, offset)
}

// GetEventsByType retrieves events of a specific type with pagination
func (s *EventService) GetEventsByType(eventType string, limit, offset int) ([]models.Event, error) {
	if eventType == "" {
		return nil, errors.New("event type is required")
	}

	if limit <= 0 {
		limit = 20 // Default limit
	}

	return s.eventRepo.GetEventsByType(eventType, limit, offset)
}

// GetRecentEvents retrieves the most recent events across all users
func (s *EventService) GetRecentEvents(limit, offset int) ([]models.Event, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}

	return s.eventRepo.GetRecentEvents(limit, offset)
}

// GetUserEventStats retrieves statistics about a user's events
func (s *EventService) GetUserEventStats(userID string) (*UserEventStats, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	stats := &UserEventStats{
		UserID:          userID,
		EventTypeCounts: make(map[string]int),
	}

	// Get total event count
	totalCount, err := s.eventRepo.CountEventsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count user events: %w", err)
	}
	stats.TotalEvents = totalCount

	// Get counts by event type
	eventTypes := []string{
		models.EventTypePageView,
		models.EventTypeProductView,
		models.EventTypeCartAdd,
		models.EventTypeCartRemove,
		models.EventTypePurchase,
		models.EventTypeWishlistAdd,
		models.EventTypeSearch,
		models.EventTypeCheckoutStart,
	}

	for _, eventType := range eventTypes {
		count, err := s.eventRepo.CountEventsByUserAndType(userID, eventType)
		if err != nil {
			continue // Skip on error
		}
		if count > 0 {
			stats.EventTypeCounts[eventType] = count
		}
	}

	return stats, nil
}

// ValidateEventType validates if the event type is supported
func (s *EventService) ValidateEventType(eventType string) error {
	supportedTypes := map[string]bool{
		models.EventTypePageView:      true,
		models.EventTypeProductView:   true,
		models.EventTypeCartAdd:       true,
		models.EventTypeCartRemove:    true,
		models.EventTypePurchase:      true,
		models.EventTypeWishlistAdd:   true,
		models.EventTypeSearch:        true,
		models.EventTypeCheckoutStart: true,
	}

	if !supportedTypes[eventType] {
		return fmt.Errorf("unsupported event type: %s", eventType)
	}

	return nil
}

// segmentsChanged checks if the user's segments have changed
func (s *EventService) segmentsChanged(oldSegments, newSegments []models.Segment) bool {
	if len(oldSegments) != len(newSegments) {
		return true
	}

	// Create maps for easier comparison
	oldMap := make(map[uuid.UUID]bool)
	for _, segment := range oldSegments {
		oldMap[segment.ID] = true
	}

	for _, segment := range newSegments {
		if !oldMap[segment.ID] {
			return true
		}
	}

	return false
}

// updateUserInterests updates user interests based on the event
func (s *EventService) updateUserInterests(event *models.Event, segments []models.Segment) {
	// Only update interests for certain event types
	if event.EventType != models.EventTypeProductView &&
		event.EventType != models.EventTypeCartAdd &&
		event.EventType != models.EventTypePurchase {
		return
	}

	// Calculate interest score based on event type
	var score float64
	switch event.EventType {
	case models.EventTypeProductView:
		score = 0.1
	case models.EventTypeCartAdd:
		score = 0.3
	case models.EventTypePurchase:
		score = 0.5
	default:
		score = 0.05
	}

	// Update interests for each segment
	for _, segment := range segments {
		if event.ProductID != nil {
			interest := models.NewUserInterest(event.UserID, segment.ID, event.ProductID, nil, score)
			s.segmentRepo.CreateUserInterest(interest)
		}

		if event.Category != nil {
			interest := models.NewUserInterest(event.UserID, segment.ID, nil, event.Category, score)
			s.segmentRepo.CreateUserInterest(interest)
		}
	}
}

// UserEventStats represents statistics about a user's events
type UserEventStats struct {
	UserID          string         `json:"user_id"`
	TotalEvents     int            `json:"total_events"`
	EventTypeCounts map[string]int `json:"event_type_counts"`
}
