package dto

import (
	"time"

	"github.com/refda/backend/internal/domain"
)

type SendOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type VerifyOTPRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Code     string `json:"code" validate:"required,len=4"`
	FullName string `json:"full_name"`
}

type AuthResponse struct {
	Token        string       `json:"token"`
	User         UserResponse `json:"user"`
	IsNewUser    bool         `json:"is_new_user"`
}

type UserResponse struct {
	ID                string `json:"id"`
	Phone             string `json:"phone"`
	FullName          string `json:"full_name"`
	Location          string `json:"location,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"omitempty,min=2,max=255"`
	Location string `json:"location" validate:"omitempty,max=255"`
}

type GiftItemRequest struct {
	Name         string   `json:"name" validate:"required,min=1,max=255"`
	Type         string   `json:"type" validate:"required,oneof=product cash"`
	ImageURLs    []string `json:"image_urls"`
	TargetAmount *float64 `json:"target_amount"`
}

type CreateEventRequest struct {
	Title        string            `json:"title" validate:"required,min=2,max=255"`
	EventType    string            `json:"event_type" validate:"required,oneof=wedding newborn birthday other"`
	EventDate    string            `json:"event_date" validate:"required"`
	EventTime    string            `json:"event_time"`
	TargetAmount *float64          `json:"target_amount"`
	Privacy      string            `json:"privacy" validate:"omitempty,oneof=public private"`
	Gifts        []GiftItemRequest `json:"gifts"`
}

type AdminReviewRequest struct {
	Approve bool `json:"approve"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type WithdrawalResponse struct {
	WithdrawalID    string  `json:"withdrawal_id"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
	EventStatus     string  `json:"event_status"`
	AmountCollected float64 `json:"amount_collected"`
	AmountWithdrawn float64 `json:"amount_withdrawn"`
	AmountAvailable float64 `json:"amount_available"`
}

type ContributeRequest struct {
	GiftID      *string `json:"gift_id"`
	EventID    string   `json:"event_id" validate:"required"`
	Amount     float64  `json:"amount" validate:"required,gt=0"`
	HideAmount bool     `json:"hide_amount"`
	Message    string   `json:"message" validate:"omitempty,max=500"`
	ReserveFull bool    `json:"reserve_full"`
}

type PayRequest struct {
	ContributionID string `json:"contribution_id" validate:"required"`
	PaymentMethod  string `json:"payment_method" validate:"required,oneof=mada apple_pay"`
}

type PaymentResponse struct {
	ReferenceNumber string                 `json:"reference_number"`
	Status          string                 `json:"status"`
	GiftCard        map[string]interface{} `json:"gift_card,omitempty"`
	Contribution    ContributionResponse   `json:"contribution"`
}

type ContributionResponse struct {
	ID              string  `json:"id"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
	ReferenceNumber string  `json:"reference_number"`
	HideAmount      bool    `json:"hide_amount"`
}

func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:                u.ID.String(),
		Phone:             u.Phone,
		FullName:          u.FullName,
		Location:          u.Location,
		ProfilePictureURL: u.ProfilePictureURL,
	}
}

func ToEventResponse(e *domain.Event, hideAmounts bool, amountCollected float64) map[string]interface{} {
	gifts := make([]map[string]interface{}, 0, len(e.Gifts))
	for _, g := range e.Gifts {
		item := map[string]interface{}{
			"id":              g.ID.String(),
			"name":            g.Name,
			"type":            g.Type,
			"image_urls":      g.ImageURLs,
			"target_amount":   g.TargetAmount,
			"status":          g.Status,
		}
		if !hideAmounts {
			item["amount_received"] = g.AmountReceived
		}
		if g.ReservedBy != nil {
			item["reserved"] = true
		}
		gifts = append(gifts, item)
	}

	resp := map[string]interface{}{
		"id":               e.ID.String(),
		"title":            e.Title,
		"event_type":       e.EventType,
		"event_date":       e.EventDate.Format(time.RFC3339),
		"event_time":       e.EventTime,
		"target_amount":    e.TargetAmount,
		"privacy":          e.Privacy,
		"status":           e.Status,
		"amount_collected": amountCollected,
		"amount_withdrawn": e.AmountWithdrawn,
		"amount_available": amountCollected - e.AmountWithdrawn,
		"hero_image_url":   e.HeroImageURL,
		"gifts":            gifts,
		"created_at":       e.CreatedAt,
	}
	if e.Beneficiary != nil {
		resp["beneficiary"] = ToUserResponse(e.Beneficiary)
	}
	return resp
}
