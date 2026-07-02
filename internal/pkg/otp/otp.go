package otp

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/refda/backend/internal/app/repositories"
	"github.com/refda/backend/internal/pkg/config"
	apperrors "github.com/refda/backend/internal/pkg/errors"
)

var saudiPhoneRegex = regexp.MustCompile(`^\+966[5][0-9]{8}$`)

type Provider struct {
	repo          *repositories.OTPRepository
	length        int
	ttl           time.Duration
	mockFixedCode string
}

func NewProvider(repo *repositories.OTPRepository, cfg *config.Config) *Provider {
	return &Provider{
		repo:          repo,
		length:        cfg.OTP.Length,
		ttl:           time.Duration(cfg.OTP.ExpiryMinutes) * time.Minute,
		mockFixedCode: cfg.OTP.MockFixedCode,
	}
}

func NormalizePhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	if strings.HasPrefix(phone, "05") {
		phone = "+966" + phone[1:]
	} else if strings.HasPrefix(phone, "5") && len(phone) == 9 {
		phone = "+966" + phone
	} else if strings.HasPrefix(phone, "966") {
		phone = "+" + phone
	}

	if !saudiPhoneRegex.MatchString(phone) {
		return "", apperrors.ErrInvalidPhone
	}
	return phone, nil
}

func (p *Provider) Send(ctx context.Context, phone string) (string, error) {
	code := p.generate()
	expiresAt := time.Now().Add(p.ttl)

	if err := p.repo.Upsert(ctx, phone, code, expiresAt); err != nil {
		return "", err
	}

	log.Printf("[MOCK SMS] OTP for %s: %s (expires in %s)", phone, code, p.ttl)
	return code, nil
}

func (p *Provider) Verify(ctx context.Context, phone, code string) error {
	challenge, err := p.repo.FindValid(ctx, phone)
	if err != nil {
		return err
	}
	if challenge == nil {
		return apperrors.ErrOTPExpired
	}
	if challenge.Code != code {
		return apperrors.ErrInvalidOTP
	}
	return p.repo.Delete(ctx, phone)
}

func (p *Provider) generate() string {
	if p.mockFixedCode != "" {
		return p.mockFixedCode
	}
	max := 1
	for i := 0; i < p.length; i++ {
		max *= 10
	}
	n := rand.Intn(max)
	return fmt.Sprintf("%0*d", p.length, n)
}
