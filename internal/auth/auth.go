package auth

import (
	"time"
)

type Auth struct {
	hmacTokenSecret []byte
	tokenAge        time.Duration
}

type Claims struct {
	UserId string
	Exp    time.Time
	Iat    time.Time
}

func NewAuth(hmacTokenSecret []byte, tokenAge time.Duration) *Auth {
	return &Auth{
		hmacTokenSecret: hmacTokenSecret,
		tokenAge:        tokenAge,
	}
}
