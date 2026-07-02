package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Phone             string    `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	FullName          string    `gorm:"size:255;not null" json:"full_name"`
	Location          string    `gorm:"size:255" json:"location,omitempty"`
	ProfilePictureURL string    `gorm:"size:512" json:"profile_picture_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type EventType string

const (
	EventTypeWedding  EventType = "wedding"
	EventTypeNewborn  EventType = "newborn"
	EventTypeBirthday EventType = "birthday"
	EventTypeOther    EventType = "other"
)

type Privacy string

const (
	PrivacyPublic  Privacy = "public"
	PrivacyPrivate Privacy = "private"
)

type EventStatus string

const (
	EventStatusPendingReview    EventStatus = "pending_review"
	EventStatusOpen             EventStatus = "open"
	EventStatusCollected        EventStatus = "collected"
	EventStatusPartialWithdrawn EventStatus = "partial_withdrawn"
	EventStatusFullyWithdrawn   EventStatus = "fully_withdrawn"
	EventStatusRejected         EventStatus = "rejected"
)

type Event struct {
	ID               uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID   `gorm:"type:uuid;index;not null" json:"user_id"`
	Title            string      `gorm:"size:255;not null" json:"title"`
	EventType        EventType   `gorm:"size:32;not null" json:"event_type"`
	EventDate        time.Time   `gorm:"not null" json:"event_date"`
	EventTime        string      `gorm:"size:8" json:"event_time,omitempty"`
	TargetAmount     *float64    `json:"target_amount,omitempty"`
	Privacy          Privacy     `gorm:"size:16;not null;default:public" json:"privacy"`
	Status           EventStatus `gorm:"size:32;not null;default:pending_review" json:"status"`
	AmountWithdrawn  float64     `gorm:"default:0;not null" json:"amount_withdrawn"`
	HeroImageURL     string      `gorm:"size:512" json:"hero_image_url,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Gifts            []Gift      `gorm:"foreignKey:EventID" json:"gifts,omitempty"`
	Beneficiary      *User       `gorm:"foreignKey:UserID" json:"beneficiary,omitempty"`
	WithdrawalRequests []WithdrawalRequest `gorm:"foreignKey:EventID" json:"-"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type GiftType string

const (
	GiftTypeProduct GiftType = "product"
	GiftTypeCash    GiftType = "cash"
)

type GiftStatus string

const (
	GiftStatusAvailable  GiftStatus = "available"
	GiftStatusReserved   GiftStatus = "reserved"
	GiftStatusFunded     GiftStatus = "funded"
	GiftStatusPartial    GiftStatus = "partial"
)

type Gift struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	EventID        uuid.UUID  `gorm:"type:uuid;index;not null" json:"event_id"`
	Name           string     `gorm:"size:255;not null" json:"name"`
	Type           GiftType   `gorm:"size:16;not null" json:"type"`
	ImageURLs      string     `gorm:"type:text" json:"image_urls,omitempty"`
	TargetAmount   *float64   `json:"target_amount,omitempty"`
	AmountReceived float64    `gorm:"default:0" json:"amount_received"`
	ReservedBy     *uuid.UUID `gorm:"type:uuid" json:"reserved_by,omitempty"`
	Status         GiftStatus `gorm:"size:16;default:available" json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (g *Gift) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

type ContributionStatus string

const (
	ContributionStatusPending   ContributionStatus = "pending"
	ContributionStatusCompleted ContributionStatus = "completed"
	ContributionStatusFailed    ContributionStatus = "failed"
)

type PaymentMethod string

const (
	PaymentMethodMada      PaymentMethod = "mada"
	PaymentMethodApplePay  PaymentMethod = "apple_pay"
	PaymentMethodCash      PaymentMethod = "cash"
)

type Contribution struct {
	ID              uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	GiftID          *uuid.UUID         `gorm:"type:uuid;index" json:"gift_id,omitempty"`
	EventID         uuid.UUID          `gorm:"type:uuid;index;not null" json:"event_id"`
	GifterID        *uuid.UUID         `gorm:"type:uuid;index" json:"gifter_id,omitempty"`
	Amount          float64            `gorm:"not null" json:"amount"`
	ReferenceNumber string             `gorm:"uniqueIndex;size:32;not null" json:"reference_number"`
	Status          ContributionStatus `gorm:"size:16;default:pending" json:"status"`
	PaymentMethod   PaymentMethod      `gorm:"size:16" json:"payment_method,omitempty"`
	HideAmount      bool               `gorm:"default:false" json:"hide_amount"`
	Message         string             `gorm:"size:500" json:"message,omitempty"`
	GiftCardMeta    string             `gorm:"type:text" json:"gift_card_meta,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Gift            *Gift              `gorm:"foreignKey:GiftID" json:"gift,omitempty"`
}

func (c *Contribution) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type MediaEntityType string

const (
	MediaEntityUser  MediaEntityType = "user"
	MediaEntityEvent MediaEntityType = "event"
	MediaEntityGift  MediaEntityType = "gift"
)

type Media struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	EntityType MediaEntityType `gorm:"size:16;not null" json:"entity_type"`
	EntityID   uuid.UUID       `gorm:"type:uuid;index;not null" json:"entity_id"`
	URL        string          `gorm:"size:512;not null" json:"url"`
	CreatedAt  time.Time       `json:"created_at"`
}

func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type WithdrawalStatus string

const (
	WithdrawalStatusPending   WithdrawalStatus = "pending"
	WithdrawalStatusCompleted WithdrawalStatus = "completed"
	WithdrawalStatusFailed    WithdrawalStatus = "failed"
)

type WithdrawalRequest struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	EventID   uuid.UUID        `gorm:"type:uuid;index;not null" json:"event_id"`
	UserID    uuid.UUID        `gorm:"type:uuid;index;not null" json:"user_id"`
	Amount    float64          `gorm:"not null" json:"amount"`
	Status    WithdrawalStatus `gorm:"size:16;not null;default:pending" json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (w *WithdrawalRequest) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
