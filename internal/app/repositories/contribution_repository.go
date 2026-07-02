package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/domain"
	"gorm.io/gorm"
)

type ContributionRepository struct {
	db *gorm.DB
}

func NewContributionRepository(db *gorm.DB) *ContributionRepository {
	return &ContributionRepository{db: db}
}

func (r *ContributionRepository) Create(ctx context.Context, c *domain.Contribution) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ContributionRepository) Update(ctx context.Context, c *domain.Contribution) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *ContributionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Contribution, error) {
	var c domain.Contribution
	err := r.db.WithContext(ctx).Preload("Gift").First(&c, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &c, err
}

func (r *ContributionRepository) FindByReference(ctx context.Context, ref string) (*domain.Contribution, error) {
	var c domain.Contribution
	err := r.db.WithContext(ctx).Preload("Gift").Where("reference_number = ?", ref).First(&c).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &c, err
}

func (r *ContributionRepository) SumCompletedByEvent(ctx context.Context, eventID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&domain.Contribution{}).
		Where("event_id = ? AND status = ?", eventID, domain.ContributionStatusCompleted).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *ContributionRepository) DB() *gorm.DB {
	return r.db
}
