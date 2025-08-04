package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mitjabez/bite-tracker/internal/db"
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
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeout)
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
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeout)
	defer cancel()
	params := sqlc.CreateUserParams{
		Email:        email,
		FullName:     fullName,
		PasswordHash: passwordHash,
	}
	user, err := r.dbContext.Queries.CreateUser(ctx, params)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Id:           user.ID,
		FullName:     user.FullName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, userId uuid.UUID, fullName string, email string, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, db.WriteTimeout)
	defer cancel()
	params := sqlc.UpdateUserParams{
		ID:           userId,
		Email:        email,
		FullName:     fullName,
		PasswordHash: passwordHash,
	}
	return r.dbContext.Queries.UpdateUser(ctx, params)
}

func (r *UserRepo) GetUser(ctx context.Context, userId uuid.UUID) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeout)
	defer cancel()
	user, err := r.dbContext.Queries.GetUser(ctx, userId)
	if err == pgx.ErrNoRows {
		return model.User{}, ErrNotFound
	} else if err != nil {
		return model.User{}, err
	}

	return mapUser(user), nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeout)
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
	return model.User{
		Id:           user.ID,
		FullName:     user.FullName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}
