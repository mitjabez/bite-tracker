package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitjabez/bite-tracker/internal/config"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
)

const maxTries = 3

type DBContext struct {
	Queries *sqlc.Queries
	Pool    *pgxpool.Pool
}

func Init(config config.Config) (DBContext, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return DBContext{}, err
	}

	for i := 1; i <= maxTries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), WriteTimeout)
		err = pool.Ping(ctx)
		cancel()

		if err != nil {
			log.Printf("DB ping failed (%d/%d): %v\n", i, maxTries, err)
		} else {
			break
		}

		if i < 3 {
			time.Sleep(time.Duration(i*2) * time.Second)
		}
	}

	if err != nil {
		return DBContext{}, err
	}

	return DBContext{
		Queries: sqlc.New(pool),
		Pool:    pool,
	}, nil
}

func RunMigration(config config.Config) error {
	// Migration is run with dedicated account
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBMigrationUsername,
		config.DBMigrationPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	m, err := migrate.New(
		"file://internal/db/migrations",
		connectionString)
	if err != nil {
		return err
	}

	log.Printf("Running DB migration ...\n")
	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Printf("No DB migration needed\n")
		return nil
	} else if err != nil {
		return err
	}
	log.Printf("DB migration successful!")
	return nil
}
