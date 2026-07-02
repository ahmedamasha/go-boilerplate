package repositories

import (
	"context"
	"time"

	"github.com/refda/backend/internal/domain"
	"gorm.io/gorm"
)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Upsert(ctx context.Context, phone, code string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("phone = ?", phone).Delete(&domain.OTPChallenge{}).Error; err != nil {
			return err
		}
		return tx.Create(&domain.OTPChallenge{
			Phone:     phone,
			Code:      code,
			ExpiresAt: expiresAt,
		}).Error
	})
}

func (r *OTPRepository) FindValid(ctx context.Context, phone string) (*domain.OTPChallenge, error) {
	var challenge domain.OTPChallenge
	err := r.db.WithContext(ctx).
		Where("phone = ? AND expires_at > ?", phone, time.Now()).
		First(&challenge).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &challenge, err
}

func (r *OTPRepository) Delete(ctx context.Context, phone string) error {
	return r.db.WithContext(ctx).Where("phone = ?", phone).Delete(&domain.OTPChallenge{}).Error
}
