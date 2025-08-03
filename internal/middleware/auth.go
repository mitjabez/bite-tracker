package middleware

import (
	"log"
	"net/http"

	"github.com/mitjabez/bite-tracker/internal/auth"
	"github.com/mitjabez/bite-tracker/internal/model"
)

func (m *Middleware) authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.auth.VerifyToken(r)
		if err == http.ErrNoCookie {
			log.Println("No auth cookie found")
			loginRedirect(w, r)
			return
		} else if err == auth.ErrTokenExpired {
			log.Println("Token expired: ", claims.Exp)
			loginRedirect(w, r)
		} else if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		user := model.User{
			Id:       claims.UserId,
			FullName: claims.FullName,
			Email:    claims.Email,
		}
		ctx := m.auth.PutUserToContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loginRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", 302)
}
