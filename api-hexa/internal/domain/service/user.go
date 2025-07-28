package service

import (
	"context"
	"errors"

	"github.com/tasjen/fz/api-hexa/internal/domain/model"
	"github.com/tasjen/fz/api-hexa/internal/domain/port"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, user model.User) error {
	if err := validateUser(user); err != nil {
		return err
	}
	_, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func validateUser(user model.User) error {
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(user.Username) < 3 || len(user.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	return nil
}
