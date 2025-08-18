package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
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
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(config.DbAppUserUsername),
		url.QueryEscape(config.DbAppUserPassword),
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbSslMode,
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
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(config.DbMigrateUserUsername),
		url.QueryEscape(config.DbMigrateUserPassword),
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbSslMode,
	)

	if config.DbBootstrapRoles {
		log.Printf("Bootstrapping DB roles\n")
		err := bootstrapDBRoles(connectionString, config)
		if err != nil {
			return err
		}
	}

	// Run standard migrations
	m, err := migrate.New("file://internal/db/migrations", connectionString)
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
	log.Printf("Successfully performed DB migration to version %d, dirty=%t\n", version, dirty)
	return nil
}

func bootstrapDBRoles(connectionString string, config config.Config) error {
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return fmt.Errorf("Cannot connect to DB for role bootstrapping: %v", err)
	}
	defer conn.Close(context.Background())

	// Add DDL (migrations) and DML permission to appuser
	// On live deployments we would use different roles for this
	username := config.DbAppUserUsername
	password := config.DbAppUserPassword
	dbName := config.DbName
	schema := "public"

	sql := fmt.Sprintf(`
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '%s') THEN
        CREATE ROLE %s LOGIN PASSWORD '%s';
    ELSE
        ALTER ROLE %s WITH PASSWORD '%s';
    END IF;

    GRANT CONNECT ON DATABASE %s TO %s;

    GRANT USAGE, CREATE ON SCHEMA %s TO %s;

    -- basic DML
    GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA %s TO %s;
    ALTER DEFAULT PRIVILEGES IN SCHEMA %s
        GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %s;

    -- migrations (DDL)
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA %s TO %s;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA %s TO %s;
    GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA %s TO %s;

    ALTER DEFAULT PRIVILEGES IN SCHEMA %s
        GRANT ALL PRIVILEGES ON TABLES TO %s;
    ALTER DEFAULT PRIVILEGES IN SCHEMA %s
        GRANT ALL PRIVILEGES ON SEQUENCES TO %s;
    ALTER DEFAULT PRIVILEGES IN SCHEMA %s
        GRANT ALL PRIVILEGES ON FUNCTIONS TO %s;
END
$$;`,
		username,
		username, password,
		username, password,
		dbName, username,
		schema, username,
		schema, username,
		schema, username,
		schema, username,
		schema, username,
		schema, username,
		schema, username,
		schema, username,
		schema, username,
	)

	ctx, close := context.WithTimeout(context.Background(), WriteTimeout)
	defer close()
	if _, err := conn.Exec(ctx, sql); err != nil {
		return fmt.Errorf("Failed to bootstrap role %s: %w", username, err)
	}
	return nil
}
