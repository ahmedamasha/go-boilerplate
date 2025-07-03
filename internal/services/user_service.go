package services

import (
	"cusror_ai/internal/models"
	"cusror_ai/internal/repositories"
	"database/sql"
	"errors"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("name and email and password are required")
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserService) UpdateUser(user *models.User) (*models.User, error) {
	if user.Name == "" || user.Email == "" {
		return nil, errors.New("name and email are required")
	}

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(id int) error {
	err := s.userRepo.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}
