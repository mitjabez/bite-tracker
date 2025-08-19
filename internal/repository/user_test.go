package repository

import (
	"context"
	"log"
	"path/filepath"
	"testing"

	"github.com/mitjabez/bite-tracker/internal/config"
	"github.com/mitjabez/bite-tracker/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestContainer(t *testing.T) {
	ctx := context.Background()

	dbName := "users"
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
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	dbPort, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)
	dbHost, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	repo := NewUserRepo(&dbContext)
	user, err := repo.GetUserByEmail(ctx, "sj@dot.com")
	assert.NoError(t, err)
	assert.Equal(t, "Salsa Jimmy", user.FullName)
}
