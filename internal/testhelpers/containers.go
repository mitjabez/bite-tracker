package testhelpers

import (
	"context"
	"path/filepath"

	"github.com/mitjabez/bite-tracker/internal/config"
	"github.com/mitjabez/bite-tracker/internal/db"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContext struct {
	*postgres.PostgresContainer
	DBContext db.DBContext
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContext, error) {
	dbName := "bite_tracker"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-users-db.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	// defer func() {
	// 	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
	// 		log.Printf("failed to terminate container: %s", err)
	// 	}
	// }()
	if err != nil {
		return nil, err
	}

	dbPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}
	dbHost, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	config := config.Config{
		ListenAddr:            "",
		HmacTokenSecret:       "",
		TokenAge:              0,
		DbName:                dbName,
		DbHost:                dbHost,
		DbPort:                dbPort.Int(),
		DbSslMode:             "disable",
		DbAppUserUsername:     dbUser,
		DbAppUserPassword:     dbPassword,
		DbMigrateUserUsername: dbUser,
		DbMigrateUserPassword: dbPassword,
		DbBootstrapRoles:      false,
	}

	dbContext, err := db.Init(config)
	if err != nil {
		return nil, err
	}
	return &PostgresContext{
		PostgresContainer: postgresContainer,
		DBContext:         dbContext,
	}, nil
}
