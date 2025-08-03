package middleware

import (
	"log"
	"net/http"

	"github.com/mitjabez/bite-tracker/internal/auth"
	"github.com/mitjabez/bite-tracker/internal/httpx"
	"github.com/mitjabez/bite-tracker/internal/model"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, model.User)

func (m *Middleware) authHandler(next AuthenticatedHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.auth.VerifyToken(r)
		if err == http.ErrNoCookie {
			log.Println("No auth cookie found")
			loginRedirect(w, r)
			return
		} else if err == auth.ErrTokenExpired {
			log.Println("Token expired: ", claims.Exp)
			m.auth.InvalidateCookieToken(w)
			loginRedirect(w, r)
			return
		} else if err != nil {
			httpx.InternalError(w, "Cannot verify token", err)
			return
		}

		user := model.User{
			Id:       claims.UserId,
			FullName: claims.FullName,
			Email:    claims.Email,
		}
		next(w, r, user)
	})
}

func loginRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", 302)
}
