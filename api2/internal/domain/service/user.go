package service

import (
	"context"
	"errors"

	"github.com/tasjen/fz/api-hexa/internal/domain/model"
	"github.com/tasjen/fz/api-hexa/internal/domain/port"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, user model.User) (model.User, error) {
	if err := validateUser(user); err != nil {
		return model.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	hashedPasswordStr := string(hashedPassword)
	user.Password = &hashedPasswordStr

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return createdUser, nil
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
