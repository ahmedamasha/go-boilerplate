package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPChallenge struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Phone     string    `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Code      string    `gorm:"size:10;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (o *OTPChallenge) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (o *OTPChallenge) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}
