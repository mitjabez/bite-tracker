package middleware

import (
	"net/http"

	"github.com/mitjabez/bite-tracker/internal/auth"
)

var middlewares = []func(next http.Handler) http.Handler{
	// Middlewares execute in ascending order
	logger,
	errorRecovery,
}

type Middleware struct {
	isAuthentication bool
	auth             *auth.Auth
}

func NewChainWithAuth(auth *auth.Auth) Middleware {
	return Middleware{
		isAuthentication: true,
		auth:             auth,
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

	if m.isAuthentication {
		return m.authHandler(n)
	}

	return n
}
