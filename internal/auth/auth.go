package auth

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	hmacTokenSecret []byte
	tokenAge        time.Duration
}

type Claims struct {
	UserId   uuid.UUID
	FullName string
	Email    string
	Exp      time.Time
	Iat      time.Time
}

func NewAuth(hmacTokenSecret string, tokenAge time.Duration) *Auth {
	return &Auth{
		hmacTokenSecret: []byte(hmacTokenSecret),
		tokenAge:        tokenAge,
	}
}
