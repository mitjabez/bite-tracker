package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrTokenExpired = errors.New("token expired")

func (a *Auth) VerifyToken(r *http.Request) (Claims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return Claims{}, http.ErrNoCookie
		}
		return Claims{}, err
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
		return a.hmacTokenSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, errors.New("invalid Claims format")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return Claims{}, fmt.Errorf("missing or invalid sub: %v", claims["sub"])
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return Claims{}, errors.New("missing or invalid exp")
	}
	iat, ok := claims["iat"].(float64)
	if !ok {
		return Claims{}, errors.New("missing or invalid iat")
	}

	userId, err := uuid.Parse(sub)
	if err != nil {
		return Claims{}, fmt.Errorf("sub claim is not in correct format" + sub)
	}
	tokenClaims := Claims{
		UserId: userId,
		Exp:    time.Unix(int64(exp), 0),
		Iat:    time.Unix(int64(iat), 0),
	}

	if tokenClaims.Exp.Before(time.Now()) {
		return tokenClaims, ErrTokenExpired
	}

	return tokenClaims, nil
}

func (a *Auth) IssueCookieToken(userId uuid.UUID) (http.Cookie, error) {
	now := time.Now()
	exp := now.Add(time.Duration(a.tokenAge))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	})

	tokenString, err := token.SignedString(a.hmacTokenSecret)
	if err != nil {
		return http.Cookie{}, err
	}

	return http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Path:    "/",
		Expires: exp,
		MaxAge:  int(a.tokenAge / time.Second),
		// TODO: Enforce true for prod
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}, nil
}

func (a *Auth) InvalidateCookieToken() http.Cookie {
	return http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  0,
		// TODO: Enforce true for prod
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
