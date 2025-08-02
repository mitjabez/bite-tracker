package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
	"github.com/mitjabez/bite-tracker/internal/model"
)

type UserRepo struct {
	dbContext *db.DBContext
}

func NewUserRepo(dbContext *db.DBContext) *UserRepo {
	return &UserRepo{dbContext}
}

func (r *UserRepo) UserExists(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_, err := r.dbContext.Queries.GetUserByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, fullName string, email string, passwordHash string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	params := sqlc.CreateUserParams{
		Email:        email,
		FullName:     fullName,
		PasswordHash: &passwordHash,
	}
	user, err := r.dbContext.Queries.CreateUser(ctx, params)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Id:           user.ID.String(),
		FullName:     user.FullName,
		Email:        user.Email,
		PasswordHash: *user.PasswordHash,
	}, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, userId string, fullName string, email string, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	params := sqlc.UpdateUserParams{
		ID:           uuid.MustParse(userId),
		Email:        email,
		FullName:     fullName,
		PasswordHash: &passwordHash,
	}
	return r.dbContext.Queries.UpdateUser(ctx, params)
}

func (r *UserRepo) GetUser(ctx context.Context, userId string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	user, err := r.dbContext.Queries.GetUser(ctx, uuid.MustParse(userId))
	if err == pgx.ErrNoRows {
		return model.User{}, ErrNotFound
	} else if err != nil {
		return model.User{}, err
	}

	return mapUser(user), nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	user, err := r.dbContext.Queries.GetUserByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return model.User{}, ErrNotFound
	} else if err != nil {
		return model.User{}, err
	}

	return mapUser(user), nil
}

func mapUser(user sqlc.User) model.User {
	var hash string
	if user.PasswordHash != nil {
		hash = *user.PasswordHash
	}

	return model.User{
		Id:           user.ID.String(),
		FullName:     user.FullName,
		Email:        user.Email,
		PasswordHash: hash,
	}
}
