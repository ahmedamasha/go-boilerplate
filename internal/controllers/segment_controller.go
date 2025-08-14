package controllers

import (
	"cusror_ai/internal/models"
	"cusror_ai/internal/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SegmentController struct {
	segmentService *services.SegmentService
}

func NewSegmentController(segmentService *services.SegmentService) *SegmentController {
	return &SegmentController{
		segmentService: segmentService,
	}
}

// CreateSegment handles POST /api/v1/segments
func (c *SegmentController) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var req CreateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "segment name is required", http.StatusBadRequest)
		return
	}

	if req.Rules == nil {
		http.Error(w, "segment rules are required", http.StatusBadRequest)
		return
	}

	// Create segment
	segment, err := c.segmentService.CreateSegment(req.Name, req.Description, req.Rules)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(segment)
}

// GetSegment handles GET /api/v1/segments/{id}
func (c *SegmentController) GetSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid segment ID", http.StatusBadRequest)
		return
	}

	segment, err := c.segmentService.GetSegmentByID(id)
	if err != nil {
		if err.Error() == "segment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(segment)
}

// GetAllSegments handles GET /api/v1/segments
func (c *SegmentController) GetAllSegments(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	var segments []models.Segment
	var err error

	if activeOnly {
		segments, err = c.segmentService.GetActiveSegments()
	} else {
		segments, err = c.segmentService.GetAllSegments()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"segments":    segments,
		"count":       len(segments),
		"active_only": activeOnly,
	})
}

// UpdateSegment handles PUT /api/v1/segments/{id}
func (c *SegmentController) UpdateSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid segment ID", http.StatusBadRequest)
		return
	}

	var req UpdateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing segment
	segment, err := c.segmentService.GetSegmentByID(id)
	if err != nil {
		if err.Error() == "segment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Update fields if provided
	if req.Name != nil {
		segment.Name = *req.Name
	}
	if req.Description != nil {
		segment.Description = *req.Description
	}
	if req.Rules != nil {
		segment.Rules = req.Rules
	}
	if req.IsActive != nil {
		segment.IsActive = *req.IsActive
	}

	// Update segment
	updatedSegment, err := c.segmentService.UpdateSegment(segment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSegment)
}

// DeleteSegment handles DELETE /api/v1/segments/{id}
func (c *SegmentController) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid segment ID", http.StatusBadRequest)
		return
	}

	err = c.segmentService.DeleteSegment(id)
	if err != nil {
		if err.Error() == "segment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AnalyzeUserSegments handles POST /api/v1/segments/analyze/{user_id}
func (c *SegmentController) AnalyzeUserSegments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	segments, err := c.segmentService.AnalyzeUserSegments(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  userID,
		"segments": segments,
		"count":    len(segments),
	})
}

// GetUserSegments handles GET /api/v1/segments/user/{user_id}
func (c *SegmentController) GetUserSegments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	segments, err := c.segmentService.GetUserSegments(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  userID,
		"segments": segments,
		"count":    len(segments),
	})
}

// GetSegmentUsers handles GET /api/v1/segments/{id}/users
func (c *SegmentController) GetSegmentUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid segment ID", http.StatusBadRequest)
		return
	}

	users, err := c.segmentService.GetSegmentUsers(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"segment_id": id,
		"users":      users,
		"count":      len(users),
	})
}

// GetUserInterests handles GET /api/v1/segments/user/{user_id}/interests
func (c *SegmentController) GetUserInterests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	interests, err := c.segmentService.GetUserInterests(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   userID,
		"interests": interests,
		"count":     len(interests),
	})
}

// CreatePredefinedSegments handles POST /api/v1/segments/predefined
func (c *SegmentController) CreatePredefinedSegments(w http.ResponseWriter, r *http.Request) {
	err := c.segmentService.CreatePredefinedSegments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Predefined segments created successfully",
	})
}

// GetSegmentRulesTemplate handles GET /api/v1/segments/rules/template
func (c *SegmentController) GetSegmentRulesTemplate(w http.ResponseWriter, r *http.Request) {
	template := map[string]interface{}{
		"rules": []map[string]interface{}{
			{
				"field":       "total_spent",
				"operator":    "gte",
				"value":       500.0,
				"description": "User has spent at least $500",
			},
			{
				"field":       "conversion_rate",
				"operator":    "lt",
				"value":       0.1,
				"description": "User has conversion rate less than 10%",
			},
		},
		"available_fields": []string{
			"total_events",
			"total_spent",
			"total_sessions",
			"page_views",
			"product_views",
			"cart_adds",
			"purchases",
			"avg_session_value",
			"conversion_rate",
		},
		"available_operators": []string{
			"gt",  // greater than
			"gte", // greater than or equal
			"lt",  // less than
			"lte", // less than or equal
			"eq",  // equal
			"ne",  // not equal
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// Request/Response structures

type CreateSegmentRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Rules       models.SegmentRules `json:"rules"`
}

type UpdateSegmentRequest struct {
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Rules       models.SegmentRules `json:"rules,omitempty"`
	IsActive    *bool               `json:"is_active,omitempty"`
}
