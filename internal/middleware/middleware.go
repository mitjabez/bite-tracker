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
	auth *auth.Auth
}

func New(auth *auth.Auth) Middleware {
	return Middleware{
		auth: auth,
	}
}

func (m *Middleware) AuthChain(next AuthenticatedHandler) http.Handler {
	n := m.authHandler(next)
	return m.Chain(n.ServeHTTP)
}

func (m *Middleware) Chain(next http.HandlerFunc) http.Handler {
	n := http.Handler(next)
	for i := len(middlewares) - 1; i >= 0; i-- {
		n = middlewares[i](n)
	}

	return n
}
