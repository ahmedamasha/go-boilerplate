package services

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/app/dto"
	"github.com/refda/backend/internal/app/repositories"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/refda/backend/internal/pkg/storage"
)

type UserService struct {
	users   *repositories.UserRepository
	storage storage.Storage
}

func NewUserService(users *repositories.UserRepository, store storage.Storage) *UserService {
	return &UserService{users: users, storage: store}
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}
	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest, avatar *multipart.FileHeader) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Location != "" {
		user.Location = req.Location
	}
	if avatar != nil {
		url, err := s.storage.Save(avatar, "profiles")
		if err != nil {
			return nil, err
		}
		user.ProfilePictureURL = url
	}

	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	resp := dto.ToUserResponse(user)
	return &resp, nil
}
