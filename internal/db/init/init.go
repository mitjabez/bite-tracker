package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitjabez/bite-tracker/internal/config"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
)

type DBContext struct {
	Queries *sqlc.Queries
	Pool    *pgxpool.Pool
}

func Init(config config.Config) (DBContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return DBContext{}, err
	}

	return DBContext{
		Queries: sqlc.New(pool),
		Pool:    pool,
	}, nil
}
