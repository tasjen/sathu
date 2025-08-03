package repository

import (
	"context"
	"errors"

	"github.com/tasjen/fz/api-hexa/internal/adapter/postgres"
	sqlc_gen "github.com/tasjen/fz/api-hexa/internal/adapter/postgres/sqlc/gen"
	"github.com/tasjen/fz/api-hexa/internal/domain/model"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	param := sqlc_gen.CreateUserWithPasswordParams{
		Email:        user.Email,
		PasswordHash: user.Password,
	}
	id, err := repo.db.Queries.CreateUserWithPassword(ctx, param)
	if err != nil {
		pgErr, ok := postgres.ToPgError(err)
		if ok && pgErr.ConstraintName == "email_unique" {
			return model.User{}, errors.New("email already exists")
		}
		return model.User{}, err
	}

	user.ID = id
	return user, nil
}
