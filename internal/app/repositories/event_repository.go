package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *domain.Event, gifts []domain.Gift) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(event).Error; err != nil {
			return err
		}
		for i := range gifts {
			gifts[i].EventID = event.ID
		}
		if len(gifts) > 0 {
			if err := tx.Create(&gifts).Error; err != nil {
				return err
			}
		}
		event.Gifts = gifts
		return nil
	})
}

func (r *EventRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).
		Preload("Gifts").
		Preload("Beneficiary").
		First(&event, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &event, err
}

func (r *EventRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).
		Preload("Gifts").
		Where("user_id = ?", userID).
		Order("event_date DESC").
		Find(&events).Error
	return events, err
}

func (r *EventRepository) ListPublic(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).
		Preload("Gifts").
		Preload("Beneficiary").
		Where("privacy = ? AND status NOT IN ?", domain.PrivacyPublic, []domain.EventStatus{
			domain.EventStatusPendingReview,
			domain.EventStatusRejected,
		}).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

func (r *EventRepository) Update(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *EventRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	err := tx.WithContext(ctx).
		Preload("Gifts").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&event, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &event, err
}

func (r *EventRepository) DB() *gorm.DB {
	return r.db
}
