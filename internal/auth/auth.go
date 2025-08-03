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
	UserId uuid.UUID
	Exp    time.Time
	Iat    time.Time
}

func NewAuth(hmacTokenSecret []byte, tokenAge time.Duration) *Auth {
	return &Auth{
		hmacTokenSecret: hmacTokenSecret,
		tokenAge:        tokenAge,
	}
}
