package db

import (
	"context"
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
	pool, err := pgxpool.New(context.Background(), config.DBAppUrl)
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
	// Don't run migration with app db account
	m, err := migrate.New("file://internal/db/migrations", config.DBMigrateUrl)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Printf("No DB migration needed\n")
		return nil
	} else if err != nil {
		return err
	}
	version, dirty, _ := m.Version()
	log.Printf("Successfully performed DB migration to version %d, dirty=%t.\n", version, dirty)
	return nil
}
