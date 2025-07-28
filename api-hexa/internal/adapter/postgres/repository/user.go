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
	param := sqlc_gen.CreateUserParams{
		Username: user.Username,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
	id, err := repo.db.Queries.CreateUser(ctx, param)
	if err != nil {
		pgErr, ok := postgres.ToPgError(err)
		if ok && pgErr.Code == "23505" {
			return model.User{}, errors.New("user already exists")
		}
		return model.User{}, err
	}

	user.ID = id
	return user, nil
}
