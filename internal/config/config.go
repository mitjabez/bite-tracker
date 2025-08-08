package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr      string        `required:"true"`
	DBAppUrl        string        `required:"true"`
	DBMigrateUrl    string        `required:"true"`
	HmacTokenSecret string        `required:"true"`
	TokenAge        time.Duration `required:"true"`
}

func Init() (Config, error) {
	btConfig := Config{}
	err := envconfig.Process("bt", &btConfig)
	if err != nil {
		return Config{}, err
	}
	return btConfig, nil
}
