package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(hmacTokenSecret []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				log.Println("No auth cookie found")
				loginRedirect(w, r)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
			return hmacTokenSecret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			log.Printf("Invalid token %s\n", err)
			loginRedirect(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Error obtaining claims: ", err)
			loginRedirect(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", claims["sub"])
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func loginRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", 302)
}
