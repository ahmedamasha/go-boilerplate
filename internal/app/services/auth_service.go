package services

import (
	"context"
	"strings"

	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/repositories"
	"github.com/refda/backend/internal/domain"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	jwtpkg "github.com/refda/backend/internal/pkg/jwt"
	"github.com/refda/backend/internal/pkg/otp"
)

type AuthService struct {
	users *repositories.UserRepository
	otp   *otp.Provider
	jwt   *jwtpkg.Manager
}

func NewAuthService(users *repositories.UserRepository, otpProvider *otp.Provider, jwtMgr *jwtpkg.Manager) *AuthService {
	return &AuthService{users: users, otp: otpProvider, jwt: jwtMgr}
}

func (s *AuthService) SendOTP(ctx context.Context, phone string) error {
	normalized, err := otp.NormalizePhone(phone)
	if err != nil {
		return err
	}
	_, err = s.otp.Send(ctx, normalized)
	return err
}

func (s *AuthService) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.AuthResponse, error) {
	phone, err := otp.NormalizePhone(req.Phone)
	if err != nil {
		return nil, err
	}

	if err := s.otp.Verify(ctx, phone, req.Code); err != nil {
		return nil, err
	}

	user, err := s.users.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	isNew := false
	if user == nil {
		name := strings.TrimSpace(req.FullName)
		if name == "" {
			return nil, apperrors.New("VALIDATION_ERROR", "full_name is required for registration", "الاسم الكامل مطلوب للتسجيل", apperrors.ErrValidation)
		}
		user = &domain.User{Phone: phone, FullName: name}
		if err := s.users.Create(ctx, user); err != nil {
			return nil, err
		}
		isNew = true
	}

	token, err := s.jwt.Generate(user.ID, user.Phone)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:     token,
		User:      dto.ToUserResponse(user),
		IsNewUser: isNew,
	}, nil
}

// Register is an alias flow: send OTP then verify with full_name
func (s *AuthService) Register(ctx context.Context, phone, fullName string) error {
	if strings.TrimSpace(fullName) == "" {
		return apperrors.New("VALIDATION_ERROR", "full_name is required", "الاسم الكامل مطلوب", apperrors.ErrValidation)
	}
	normalized, err := otp.NormalizePhone(phone)
	if err != nil {
		return err
	}
	existing, err := s.users.FindByPhone(ctx, normalized)
	if err != nil {
		return err
	}
	if existing != nil {
		return apperrors.ErrUserExists
	}
	_, err = s.otp.Send(ctx, normalized)
	return err
}
