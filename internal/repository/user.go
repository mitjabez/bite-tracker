package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
)

type UserRepo struct {
	DBContext db.DBContext
}

func (r *UserRepo) UserExists(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_, err := r.DBContext.Queries.GetUser(ctx, email)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, fullName string, email string, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	params := sqlc.CreateUserParams{
		Email:        email,
		FullName:     fullName,
		PasswordHash: &passwordHash,
	}
	_, err := r.DBContext.Queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	return err
}
