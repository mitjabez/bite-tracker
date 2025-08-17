package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr            string        `required:"true" split_words:"true"`
	HmacTokenSecret       string        `required:"true" split_words:"true"`
	TokenAge              time.Duration `required:"true" split_words:"true"`
	DbName                string        `required:"true" split_words:"true"`
	DbHost                string        `required:"true" split_words:"true"`
	DbPort                int           `required:"true" split_words:"true"`
	DbSslMode             string        `required:"true" split_words:"true"`
	DbAppUserUsername     string        `required:"true" split_words:"true"`
	DbAppUserPassword     string        `required:"true" split_words:"true"`
	DbMigrateUserUsername string        `required:"true" split_words:"true"`
	DbMigrateUserPassword string        `required:"true" split_words:"true"`
	// Whether to create roles at app startup
	// Should be enabled only when starting the app for the first time or during development
	DbBootstrapRoles bool `required:"true" split_words:"true"`
}

func Init() (Config, error) {
	btConfig := Config{}
	err := envconfig.Process("bt", &btConfig)
	if err != nil {
		return Config{}, err
	}
	return btConfig, nil
}
