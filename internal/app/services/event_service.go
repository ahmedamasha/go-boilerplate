package services

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/repositories"
	"github.com/refda/backend/internal/domain"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/refda/backend/internal/pkg/storage"
	"gorm.io/gorm"
)

type EventService struct {
	events        *repositories.EventRepository
	contributions *repositories.ContributionRepository
	withdrawals   *repositories.WithdrawalRepository
	storage       storage.Storage
}

func NewEventService(
	events *repositories.EventRepository,
	contributions *repositories.ContributionRepository,
	withdrawals *repositories.WithdrawalRepository,
	store storage.Storage,
) *EventService {
	return &EventService{
		events:        events,
		contributions: contributions,
		withdrawals:   withdrawals,
		storage:       store,
	}
}

func (s *EventService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateEventRequest, heroImage *multipart.FileHeader) (map[string]interface{}, error) {
	eventDate, err := parseEventDate(req.EventDate)
	if err != nil {
		return nil, apperrors.New("VALIDATION_ERROR", "invalid event_date format (YYYY-MM-DD)", "صيغة التاريخ غير صالحة", apperrors.ErrValidation)
	}

	privacy := domain.PrivacyPublic
	if req.Privacy == string(domain.PrivacyPrivate) {
		privacy = domain.PrivacyPrivate
	}

	event := &domain.Event{
		UserID:       userID,
		Title:        req.Title,
		EventType:    domain.EventType(req.EventType),
		EventDate:    eventDate,
		EventTime:    req.EventTime,
		TargetAmount: req.TargetAmount,
		Privacy:      privacy,
		Status:       domain.EventStatusPendingReview,
	}

	if heroImage != nil {
		url, err := s.storage.Save(heroImage, "events")
		if err != nil {
			return nil, err
		}
		event.HeroImageURL = url
	}

	gifts := make([]domain.Gift, 0, len(req.Gifts))
	for _, g := range req.Gifts {
		imgJSON, _ := json.Marshal(g.ImageURLs)
		gift := domain.Gift{
			Name:         g.Name,
			Type:         domain.GiftType(g.Type),
			ImageURLs:    string(imgJSON),
			TargetAmount: g.TargetAmount,
			Status:       domain.GiftStatusAvailable,
		}
		gifts = append(gifts, gift)
	}

	if err := s.events.Create(ctx, event, gifts); err != nil {
		return nil, err
	}

	created, err := s.events.FindByID(ctx, event.ID)
	if err != nil {
		return nil, err
	}
	return s.toEventResponse(ctx, created, false)
}

func (s *EventService) GetByID(ctx context.Context, id uuid.UUID, hideAmounts bool) (map[string]interface{}, error) {
	event, err := s.events.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, apperrors.ErrEventNotFound
	}
	return s.toEventResponse(ctx, event, hideAmounts)
}

func (s *EventService) List(ctx context.Context, userID *uuid.UUID) ([]map[string]interface{}, error) {
	var events []domain.Event
	var err error

	if userID != nil {
		events, err = s.events.ListByUser(ctx, *userID)
	} else {
		events, err = s.events.ListPublic(ctx)
	}
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(events))
	for i := range events {
		item, err := s.toEventResponse(ctx, &events[i], false)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *EventService) AdminReview(ctx context.Context, eventID uuid.UUID, approve bool) (map[string]interface{}, error) {
	event, err := s.events.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, apperrors.ErrEventNotFound
	}
	if event.Status != domain.EventStatusPendingReview {
		return nil, apperrors.New("VALIDATION_ERROR", "event is not pending review", "المناسبة ليست قيد المراجعة", apperrors.ErrValidation)
	}

	if approve {
		event.Status = domain.EventStatusOpen
	} else {
		event.Status = domain.EventStatusRejected
	}

	if err := s.events.Update(ctx, event); err != nil {
		return nil, err
	}
	return s.toEventResponse(ctx, event, false)
}

func (s *EventService) RequestWithdrawal(ctx context.Context, userID, eventID uuid.UUID, amount float64) (map[string]interface{}, error) {
	if amount <= 0 {
		return nil, apperrors.New("VALIDATION_ERROR", "amount must be positive", "المبلغ يجب أن يكون موجباً", apperrors.ErrInvalidWithdrawal)
	}

	var result map[string]interface{}
	err := s.events.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		event, err := s.events.FindByIDForUpdate(ctx, tx, eventID)
		if err != nil {
			return err
		}
		if event == nil {
			return apperrors.ErrEventNotFound
		}
		if event.UserID != userID {
			return apperrors.ErrForbidden
		}
		if !CanRequestWithdrawal(event.Status) {
			return apperrors.ErrEventNotWithdrawable
		}

		collected, err := s.contributions.SumCompletedByEvent(ctx, eventID)
		if err != nil {
			return err
		}
		if !IsFullyCollected(EventTargetAmount(event), collected) {
			return apperrors.ErrEventNotWithdrawable
		}

		available := AvailableForWithdrawal(collected, event.AmountWithdrawn)
		if amount > available+moneyEpsilon {
			return apperrors.New("INVALID_WITHDRAWAL", "amount exceeds available balance", "المبلغ يتجاوز الرصيد المتاح", apperrors.ErrInvalidWithdrawal)
		}

		withdrawal := &domain.WithdrawalRequest{
			EventID: eventID,
			UserID:  userID,
			Amount:  amount,
			Status:  domain.WithdrawalStatusCompleted,
		}
		if err := tx.Create(withdrawal).Error; err != nil {
			return err
		}

		event.AmountWithdrawn += amount
		event.Status = WithdrawalEventStatus(collected, event.AmountWithdrawn)
		if err := tx.Save(event).Error; err != nil {
			return err
		}

		result = map[string]interface{}{
			"withdrawal_id":     withdrawal.ID.String(),
			"amount":            withdrawal.Amount,
			"status":            withdrawal.Status,
			"event_status":      event.Status,
			"amount_collected":  collected,
			"amount_withdrawn":  event.AmountWithdrawn,
			"amount_available":  AvailableForWithdrawal(collected, event.AmountWithdrawn),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *EventService) RecalculateAfterPayment(ctx context.Context, eventID uuid.UUID) error {
	event, err := s.events.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return apperrors.ErrEventNotFound
	}

	if event.Status != domain.EventStatusOpen {
		return nil
	}

	collected, err := s.contributions.SumCompletedByEvent(ctx, eventID)
	if err != nil {
		return err
	}

	target := EventTargetAmount(event)
	if IsFullyCollected(target, collected) {
		event.Status = domain.EventStatusCollected
		return s.events.Update(ctx, event)
	}
	return nil
}

func (s *EventService) toEventResponse(ctx context.Context, e *domain.Event, hideAmounts bool) (map[string]interface{}, error) {
	collected, err := s.contributions.SumCompletedByEvent(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	return dto.ToEventResponse(e, hideAmounts, collected), nil
}

func parseEventDate(raw string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
	}
	for _, layout := range formats {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, apperrors.ErrValidation
}
