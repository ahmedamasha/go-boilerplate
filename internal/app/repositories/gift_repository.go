package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GiftRepository struct {
	db *gorm.DB
}

func NewGiftRepository(db *gorm.DB) *GiftRepository {
	return &GiftRepository{db: db}
}

func (r *GiftRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Gift, error) {
	var gift domain.Gift
	err := r.db.WithContext(ctx).First(&gift, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &gift, err
}

func (r *GiftRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.Gift, error) {
	var gift domain.Gift
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&gift, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &gift, err
}

func (r *GiftRepository) Update(ctx context.Context, gift *domain.Gift) error {
	return r.db.WithContext(ctx).Save(gift).Error
}

func (r *GiftRepository) DB() *gorm.DB {
	return r.db
}
