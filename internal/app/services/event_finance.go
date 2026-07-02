package services

import (
	"math"

	"github.com/refda/backend/internal/domain"
)

const moneyEpsilon = 0.01

func EventTargetAmount(event *domain.Event) float64 {
	if event.TargetAmount != nil && *event.TargetAmount > 0 {
		return *event.TargetAmount
	}
	var sum float64
	for _, g := range event.Gifts {
		if g.TargetAmount != nil {
			sum += *g.TargetAmount
		}
	}
	return sum
}

func IsFullyCollected(target, collected float64) bool {
	if target <= 0 {
		return false
	}
	return collected+moneyEpsilon >= target
}

func AvailableForWithdrawal(collected, withdrawn float64) float64 {
	available := collected - withdrawn
	if available < 0 {
		return 0
	}
	return math.Round(available*100) / 100
}

func WithdrawalEventStatus(collected, withdrawn float64) domain.EventStatus {
	if collected <= 0 {
		return domain.EventStatusCollected
	}
	if withdrawn+moneyEpsilon >= collected {
		return domain.EventStatusFullyWithdrawn
	}
	return domain.EventStatusPartialWithdrawn
}

func CanAcceptContributions(status domain.EventStatus) bool {
	return status == domain.EventStatusOpen
}

func CanRequestWithdrawal(status domain.EventStatus) bool {
	return status == domain.EventStatusCollected || status == domain.EventStatusPartialWithdrawn
}
