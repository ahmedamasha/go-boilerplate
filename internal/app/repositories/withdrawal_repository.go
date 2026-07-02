package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/domain"
	"gorm.io/gorm"
)

type WithdrawalRepository struct {
	db *gorm.DB
}

func NewWithdrawalRepository(db *gorm.DB) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (r *WithdrawalRepository) Create(ctx context.Context, w *domain.WithdrawalRequest) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *WithdrawalRepository) Update(ctx context.Context, w *domain.WithdrawalRequest) error {
	return r.db.WithContext(ctx).Save(w).Error
}

func (r *WithdrawalRepository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.WithdrawalRequest, error) {
	var items []domain.WithdrawalRequest
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *WithdrawalRepository) DB() *gorm.DB {
	return r.db
}
