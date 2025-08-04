package config

import "time"

type Config struct {
	ListenAddr      string
	DBHost          string
	DBPort          int
	DBName          string
	DBUsername      string
	DBPassword      string
	HmacTokenSecret []byte
	TokenAge        time.Duration
}

func LocalDev() Config {
	return Config{
		ListenAddr:      ":8000",
		DBHost:          "localhost",
		DBPort:          5432,
		DBName:          "bite_tracker",
		DBUsername:      "biteapp",
		DBPassword:      "superburrito",
		HmacTokenSecret: []byte("1WSB6LaNNLfxi.JbTxrao0s3b4wTpH"),
		TokenAge:        time.Hour * 24,
	}
}
