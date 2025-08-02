package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

func (m *Middleware) authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.auth.VerifyToken(r)
		if err == http.ErrNoCookie {
			log.Println("No auth cookie found")
			loginRedirect(w, r)
			return
		} else if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		if claims.Exp.Before(time.Now()) {
			log.Println("Token expired: ", claims.Exp)
			loginRedirect(w, r)
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loginRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", 302)
}
