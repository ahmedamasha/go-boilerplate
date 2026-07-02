package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/repositories"
	"github.com/refda/backend/internal/domain"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"gorm.io/gorm"
)

type GiftService struct {
	gifts         *repositories.GiftRepository
	events        *repositories.EventRepository
	contributions *repositories.ContributionRepository
	eventSvc      *EventService
}

func NewGiftService(
	gifts *repositories.GiftRepository,
	events *repositories.EventRepository,
	contributions *repositories.ContributionRepository,
	eventSvc *EventService,
) *GiftService {
	return &GiftService{gifts: gifts, events: events, contributions: contributions, eventSvc: eventSvc}
}

func (s *GiftService) Contribute(ctx context.Context, gifterID *uuid.UUID, req dto.ContributeRequest) (*domain.Contribution, error) {
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, apperrors.New("VALIDATION_ERROR", "invalid event_id", "معرف المناسبة غير صالح", apperrors.ErrValidation)
	}

	event, err := s.events.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, apperrors.ErrEventNotFound
	}
	if !CanAcceptContributions(event.Status) {
		return nil, apperrors.ErrEventNotOpen
	}

	var giftID *uuid.UUID
	if req.GiftID != nil && *req.GiftID != "" {
		id, err := uuid.Parse(*req.GiftID)
		if err != nil {
			return nil, apperrors.New("VALIDATION_ERROR", "invalid gift_id", "معرف الهدية غير صالح", apperrors.ErrValidation)
		}
		giftID = &id
	}

	contribution := &domain.Contribution{
		EventID:         eventID,
		GiftID:          giftID,
		GifterID:        gifterID,
		Amount:          req.Amount,
		ReferenceNumber: generateReference(),
		Status:          domain.ContributionStatusPending,
		HideAmount:      req.HideAmount,
		Message:         req.Message,
	}

	if giftID == nil {
		if err := s.contributions.Create(ctx, contribution); err != nil {
			return nil, err
		}
		return contribution, nil
	}

	err = s.gifts.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gift, err := s.gifts.FindByIDForUpdate(ctx, tx, *giftID)
		if err != nil {
			return err
		}
		if gift == nil {
			return apperrors.ErrGiftNotFound
		}
		if gift.EventID != eventID {
			return apperrors.ErrGiftNotFound
		}

		if req.ReserveFull {
			if gift.ReservedBy != nil {
				return apperrors.ErrGiftReserved
			}
			if gift.TargetAmount != nil {
				remaining := *gift.TargetAmount - gift.AmountReceived
				if remaining <= 0 {
					return apperrors.ErrGiftReserved
				}
				contribution.Amount = remaining
			}
			gift.ReservedBy = gifterID
			gift.Status = domain.GiftStatusReserved
		} else if gift.Type == domain.GiftTypeProduct && gift.TargetAmount != nil {
			remaining := *gift.TargetAmount - gift.AmountReceived
			if req.Amount > remaining {
				return apperrors.ErrInsufficientAmount
			}
		}

		if err := tx.Save(gift).Error; err != nil {
			return err
		}
		return tx.Create(contribution).Error
	})

	if err != nil {
		return nil, err
	}
	return contribution, nil
}

func (s *GiftService) CompletePayment(ctx context.Context, contributionID uuid.UUID, method domain.PaymentMethod) (*dto.PaymentResponse, error) {
	contribution, err := s.contributions.FindByID(ctx, contributionID)
	if err != nil {
		return nil, err
	}
	if contribution == nil {
		return nil, apperrors.ErrNotFound
	}
	if contribution.Status == domain.ContributionStatusCompleted {
		return nil, apperrors.New("VALIDATION_ERROR", "contribution already paid", "تم دفع المساهمة مسبقاً", apperrors.ErrValidation)
	}

	event, err := s.events.FindByID(ctx, contribution.EventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, apperrors.ErrEventNotFound
	}
	if !CanAcceptContributions(event.Status) {
		return nil, apperrors.ErrEventNotOpen
	}

	err = s.gifts.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if contribution.GiftID != nil {
			gift, err := s.gifts.FindByIDForUpdate(ctx, tx, *contribution.GiftID)
			if err != nil {
				return err
			}
			if gift == nil {
				return apperrors.ErrGiftNotFound
			}

			if gift.TargetAmount != nil && gift.Type == domain.GiftTypeProduct {
				remaining := *gift.TargetAmount - gift.AmountReceived
				if contribution.Amount > remaining+moneyEpsilon {
					return apperrors.ErrInsufficientAmount
				}
				gift.AmountReceived += contribution.Amount
				if gift.AmountReceived+moneyEpsilon >= *gift.TargetAmount {
					gift.Status = domain.GiftStatusFunded
					gift.AmountReceived = *gift.TargetAmount
				} else {
					gift.Status = domain.GiftStatusPartial
				}
			}

			if err := tx.Save(gift).Error; err != nil {
				return err
			}
		}

		contribution.Status = domain.ContributionStatusCompleted
		contribution.PaymentMethod = method

		meta, _ := json.Marshal(map[string]interface{}{
			"title":    "هدية من رفدة",
			"title_ar": "هدية من رفدة",
			"message":  contribution.Message,
			"amount":   contribution.Amount,
			"hidden":   contribution.HideAmount,
		})
		contribution.GiftCardMeta = string(meta)

		return tx.Save(contribution).Error
	})
	if err != nil {
		return nil, err
	}

	if err := s.eventSvc.RecalculateAfterPayment(ctx, contribution.EventID); err != nil {
		return nil, err
	}

	var giftCard map[string]interface{}
	_ = json.Unmarshal([]byte(contribution.GiftCardMeta), &giftCard)

	return &dto.PaymentResponse{
		ReferenceNumber: contribution.ReferenceNumber,
		Status:          string(contribution.Status),
		GiftCard:        giftCard,
		Contribution: dto.ContributionResponse{
			ID:              contribution.ID.String(),
			Amount:          contribution.Amount,
			Status:          string(contribution.Status),
			ReferenceNumber: contribution.ReferenceNumber,
			HideAmount:      contribution.HideAmount,
		},
	}, nil
}

func generateReference() string {
	n := rand.Intn(9000) + 1000
	return fmt.Sprintf("RFD-%d", n)
}
