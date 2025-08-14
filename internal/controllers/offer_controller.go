package controllers

import (
	"cusror_ai/internal/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OfferController struct {
	offerService *services.OfferService
}

func NewOfferController(offerService *services.OfferService) *OfferController {
	return &OfferController{
		offerService: offerService,
	}
}

// GetActiveOffers handles GET /api/v1/offers
func (c *OfferController) GetActiveOffers(w http.ResponseWriter, r *http.Request) {
	offers, err := c.offerService.GetActiveOffers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"offers": offers,
		"count":  len(offers),
	})
}

// GetOffer handles GET /api/v1/offers/{id}
func (c *OfferController) GetOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid offer ID", http.StatusBadRequest)
		return
	}

	offer, err := c.offerService.GetOfferByID(id)
	if err != nil {
		if err.Error() == "offer not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offer)
}

// GetUserOffers handles GET /api/v1/offers/user/{user_id}
func (c *OfferController) GetUserOffers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	offers, err := c.offerService.GetOffersByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"offers":  offers,
		"count":   len(offers),
	})
}

// UseOffer handles POST /api/v1/offers/{id}/use
func (c *OfferController) UseOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid offer ID", http.StatusBadRequest)
		return
	}

	var req UseOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if req.Amount < 0 {
		http.Error(w, "amount must be non-negative", http.StatusBadRequest)
		return
	}

	usage, err := c.offerService.UseOffer(id, req.UserID, req.OrderID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usage)
}

// GetOfferStats handles GET /api/v1/offers/stats
func (c *OfferController) GetOfferStats(w http.ResponseWriter, r *http.Request) {
	stats, err := c.offerService.GetOfferStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Request structures
type UseOfferRequest struct {
	UserID  string  `json:"user_id"`
	OrderID *string `json:"order_id,omitempty"`
	Amount  float64 `json:"amount"`
}
