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

type SegmentService struct {
	segmentRepo *repositories.SegmentRepository
	eventRepo   *repositories.EventRepository
}

func NewSegmentService(
	segmentRepo *repositories.SegmentRepository,
	eventRepo *repositories.EventRepository,
) *SegmentService {
	return &SegmentService{
		segmentRepo: segmentRepo,
		eventRepo:   eventRepo,
	}
}

// CreateSegment creates a new segment with rules
func (s *SegmentService) CreateSegment(name, description string, rules models.SegmentRules) (*models.Segment, error) {
	if name == "" {
		return nil, errors.New("segment name is required")
	}

	// Validate rules
	if err := s.validateSegmentRules(rules); err != nil {
		return nil, fmt.Errorf("invalid rules: %w", err)
	}

	segment := models.NewSegment(name, description, rules)

	createdSegment, err := s.segmentRepo.CreateSegment(segment)
	if err != nil {
		return nil, fmt.Errorf("failed to create segment: %w", err)
	}

	return createdSegment, nil
}

// GetSegmentByID retrieves a segment by its ID
func (s *SegmentService) GetSegmentByID(id uuid.UUID) (*models.Segment, error) {
	segment, err := s.segmentRepo.GetSegmentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("segment not found")
		}
		return nil, err
	}
	return segment, nil
}

// GetAllSegments retrieves all segments
func (s *SegmentService) GetAllSegments() ([]models.Segment, error) {
	return s.segmentRepo.GetAllSegments()
}

// GetActiveSegments retrieves all active segments
func (s *SegmentService) GetActiveSegments() ([]models.Segment, error) {
	return s.segmentRepo.GetActiveSegments()
}

// UpdateSegment updates an existing segment
func (s *SegmentService) UpdateSegment(segment *models.Segment) (*models.Segment, error) {
	if segment.Name == "" {
		return nil, errors.New("segment name is required")
	}

	// Validate rules
	if err := s.validateSegmentRules(segment.Rules); err != nil {
		return nil, fmt.Errorf("invalid rules: %w", err)
	}

	segment.UpdatedAt = time.Now()

	updatedSegment, err := s.segmentRepo.UpdateSegment(segment)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("segment not found")
		}
		return nil, err
	}

	return updatedSegment, nil
}

// DeleteSegment deletes a segment
func (s *SegmentService) DeleteSegment(id uuid.UUID) error {
	err := s.segmentRepo.DeleteSegment(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("segment not found")
		}
		return err
	}
	return nil
}

// AnalyzeUserSegments analyzes which segments a user belongs to based on their behavior
func (s *SegmentService) AnalyzeUserSegments(userID string) ([]models.Segment, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get all active segments
	segments, err := s.segmentRepo.GetActiveSegments()
	if err != nil {
		return nil, fmt.Errorf("failed to get active segments: %w", err)
	}

	// Get user events for analysis
	userEvents, err := s.eventRepo.GetEventsByUserID(userID, 1000, 0) // Get last 1000 events
	if err != nil {
		return nil, fmt.Errorf("failed to get user events: %w", err)
	}

	var matchingSegments []models.Segment

	// Analyze each segment
	for _, segment := range segments {
		if s.userMatchesSegment(userID, userEvents, segment) {
			matchingSegments = append(matchingSegments, segment)

			// Assign user to segment if not already assigned
			isAssigned, _ := s.segmentRepo.IsUserInSegment(userID, segment.ID)
			if !isAssigned {
				userSegment := models.NewUserSegment(userID, segment.ID)
				s.segmentRepo.AssignUserToSegment(userSegment)
			}
		}
	}

	return matchingSegments, nil
}

// userMatchesSegment checks if a user matches the segment rules
func (s *SegmentService) userMatchesSegment(userID string, events []models.Event, segment models.Segment) bool {
	// Parse rules
	rules, ok := segment.Rules["rules"].([]interface{})
	if !ok {
		return false
	}

	// Calculate user metrics
	metrics := s.calculateUserMetrics(events)

	// Evaluate each rule
	for _, rule := range rules {
		ruleMap, ok := rule.(map[string]interface{})
		if !ok {
			continue
		}

		if s.evaluateRule(metrics, ruleMap) {
			return true
		}
	}

	return false
}

// calculateUserMetrics calculates metrics from user events
func (s *SegmentService) calculateUserMetrics(events []models.Event) map[string]interface{} {
	metrics := make(map[string]interface{})

	// Count events by type
	eventCounts := make(map[string]int)
	totalSpent := 0.0
	totalSessions := 0
	uniqueDays := make(map[string]bool)

	for _, event := range events {
		eventCounts[event.EventType]++

		// Extract day for session calculation
		day := event.Timestamp.Format("2006-01-02")
		uniqueDays[day] = true

		// Calculate total spent from purchase events
		if event.EventType == models.EventTypePurchase {
			if amount, ok := event.Metadata["amount"].(float64); ok {
				totalSpent += amount
			}
		}
	}

	totalSessions = len(uniqueDays)

	// Set metrics
	metrics["total_events"] = len(events)
	metrics["total_spent"] = totalSpent
	metrics["total_sessions"] = totalSessions
	metrics["page_views"] = eventCounts[models.EventTypePageView]
	metrics["product_views"] = eventCounts[models.EventTypeProductView]
	metrics["cart_adds"] = eventCounts[models.EventTypeCartAdd]
	metrics["purchases"] = eventCounts[models.EventTypePurchase]
	metrics["avg_session_value"] = 0.0

	if totalSessions > 0 {
		metrics["avg_session_value"] = totalSpent / float64(totalSessions)
	}

	// Calculate conversion rate
	if eventCounts[models.EventTypeProductView] > 0 {
		metrics["conversion_rate"] = float64(eventCounts[models.EventTypePurchase]) / float64(eventCounts[models.EventTypeProductView])
	} else {
		metrics["conversion_rate"] = 0.0
	}

	return metrics
}

// evaluateRule evaluates a single rule against user metrics
func (s *SegmentService) evaluateRule(metrics map[string]interface{}, rule map[string]interface{}) bool {
	field, ok := rule["field"].(string)
	if !ok {
		return false
	}

	operator, ok := rule["operator"].(string)
	if !ok {
		return false
	}

	value := rule["value"]

	userValue, exists := metrics[field]
	if !exists {
		return false
	}

	return s.compareValues(userValue, operator, value)
}

// compareValues compares two values based on the operator
func (s *SegmentService) compareValues(userValue interface{}, operator string, ruleValue interface{}) bool {
	switch operator {
	case "gt": // greater than
		return s.compareNumeric(userValue, ruleValue, func(a, b float64) bool { return a > b })
	case "gte": // greater than or equal
		return s.compareNumeric(userValue, ruleValue, func(a, b float64) bool { return a >= b })
	case "lt": // less than
		return s.compareNumeric(userValue, ruleValue, func(a, b float64) bool { return a < b })
	case "lte": // less than or equal
		return s.compareNumeric(userValue, ruleValue, func(a, b float64) bool { return a <= b })
	case "eq": // equal
		return s.compareEqual(userValue, ruleValue)
	case "ne": // not equal
		return !s.compareEqual(userValue, ruleValue)
	default:
		return false
	}
}

// compareNumeric compares numeric values
func (s *SegmentService) compareNumeric(userValue, ruleValue interface{}, compareFn func(float64, float64) bool) bool {
	uv, ok1 := s.toFloat64(userValue)
	rv, ok2 := s.toFloat64(ruleValue)

	if !ok1 || !ok2 {
		return false
	}

	return compareFn(uv, rv)
}

// compareEqual compares values for equality
func (s *SegmentService) compareEqual(userValue, ruleValue interface{}) bool {
	return userValue == ruleValue
}

// toFloat64 converts various numeric types to float64
func (s *SegmentService) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

// validateSegmentRules validates segment rules format
func (s *SegmentService) validateSegmentRules(rules models.SegmentRules) error {
	if rules == nil {
		return errors.New("rules cannot be nil")
	}

	// Check if rules contain the expected structure
	rulesList, ok := rules["rules"]
	if !ok {
		return errors.New("rules must contain 'rules' field")
	}

	// Validate rules format
	rulesArray, ok := rulesList.([]interface{})
	if !ok {
		return errors.New("rules must be an array")
	}

	for i, rule := range rulesArray {
		ruleMap, ok := rule.(map[string]interface{})
		if !ok {
			return fmt.Errorf("rule %d must be an object", i)
		}

		// Validate required fields
		if _, ok := ruleMap["field"]; !ok {
			return fmt.Errorf("rule %d missing 'field'", i)
		}
		if _, ok := ruleMap["operator"]; !ok {
			return fmt.Errorf("rule %d missing 'operator'", i)
		}
		if _, ok := ruleMap["value"]; !ok {
			return fmt.Errorf("rule %d missing 'value'", i)
		}
	}

	return nil
}

// GetUserSegments retrieves all segments for a user
func (s *SegmentService) GetUserSegments(userID string) ([]models.Segment, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.segmentRepo.GetUserSegments(userID)
}

// GetSegmentUsers retrieves all users in a segment
func (s *SegmentService) GetSegmentUsers(segmentID uuid.UUID) ([]string, error) {
	return s.segmentRepo.GetSegmentUsers(segmentID)
}

// CreatePredefinedSegments creates predefined segments with common rules
func (s *SegmentService) CreatePredefinedSegments() error {
	predefinedSegments := []struct {
		name        string
		description string
		rules       models.SegmentRules
	}{
		{
			name:        models.SegmentHighSpender,
			description: "Users who spend more than $500 on average per session",
			rules: models.SegmentRules{
				"rules": []map[string]interface{}{
					{
						"field":    "avg_session_value",
						"operator": "gte",
						"value":    500.0,
					},
				},
			},
		},
		{
			name:        models.SegmentWindowShopper,
			description: "Users who view many products but rarely purchase",
			rules: models.SegmentRules{
				"rules": []map[string]interface{}{
					{
						"field":    "product_views",
						"operator": "gte",
						"value":    10,
					},
					{
						"field":    "conversion_rate",
						"operator": "lt",
						"value":    0.1,
					},
				},
			},
		},
		{
			name:        models.SegmentDiscountSeeker,
			description: "Users who frequently search and have low conversion rates",
			rules: models.SegmentRules{
				"rules": []map[string]interface{}{
					{
						"field":    "conversion_rate",
						"operator": "lt",
						"value":    0.05,
					},
					{
						"field":    "total_sessions",
						"operator": "gte",
						"value":    5,
					},
				},
			},
		},
	}

	for _, seg := range predefinedSegments {
		// Check if segment already exists
		_, err := s.segmentRepo.GetSegmentByName(seg.name)
		if err == nil {
			continue // Segment already exists
		}

		// Create segment
		_, err = s.CreateSegment(seg.name, seg.description, seg.rules)
		if err != nil {
			return fmt.Errorf("failed to create predefined segment %s: %w", seg.name, err)
		}
	}

	return nil
}

// GetUserInterests retrieves interests for a user
func (s *SegmentService) GetUserInterests(userID string) ([]models.UserInterest, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.segmentRepo.GetUserInterests(userID)
}
