package middleware

import (
	"net/http"
)

var middlewares = []func(next http.Handler) http.Handler{
	// Middlewares execute in ascending order
	logger,
	errorRecovery,
}

type Middleware struct {
	isAuthenticated bool
	hmacTokenSecret []byte
}

func NewChainWithAuth(hmacTokenSecret []byte) Middleware {
	return Middleware{
		isAuthenticated: true,
		hmacTokenSecret: hmacTokenSecret,
	}
}

func NewChainNoAuth() Middleware {
	return Middleware{}
}

func (m *Middleware) Chain(next http.HandlerFunc) http.Handler {
	n := http.Handler(next)
	for i := len(middlewares) - 1; i >= 0; i-- {
		n = middlewares[i](n)
	}

	if m.isAuthenticated {
		return Auth(m.hmacTokenSecret, n)
	}

	return n
}
