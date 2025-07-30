package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mitjabez/bite-tracker/internal/model"
	"github.com/mitjabez/bite-tracker/internal/repository"
	"github.com/mitjabez/bite-tracker/internal/view"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo            *repository.UserRepo
	hmacTokenSecret []byte
}

func NewAuthHandler(repo *repository.UserRepo) AuthHandler {
	return AuthHandler{
		repo: repo,
		// TODO: move to config
		hmacTokenSecret: []byte("1WSB6LaNNLfxi.JbTxrao0s3b4wTpH"),
	}
}

func (h *AuthHandler) RegisterUserForm(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.RegisterUserForm(model.User{}, map[string]string{}), "Register User").Render(r.Context(), w)
}

func (h *AuthHandler) HandleRegisterUserForm(w http.ResponseWriter, r *http.Request) {
	errors := map[string]string{}

	auth := model.User{
		FullName: r.FormValue("full-name"),
		Email:    r.FormValue("email"),
	}
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	if len(auth.FullName) < 5 {
		errors["full-name"] = "Full name must be at least 5 characters long."
	}

	_, err := mail.ParseAddress(auth.Email)
	if err != nil {
		errors["email"] = "Invalid email address"
	}

	if password != confirmPassword {
		errors["password"] = "Passwords do not match"
	}

	if len(password) < 10 {
		errors["password"] = "Password must be at least 10 characters long."
	}

	if len(errors) == 0 {
		userExists, err := h.repo.UserExists(r.Context(), auth.Email)
		if err != nil {
			log.Fatal("Error checking if user exists: ", err)
		}
		if userExists {
			errors["email"] = "Username unavailable — please choose another"
		}
	}

	if len(errors) > 0 {
		view.Layout(view.RegisterUserForm(auth, errors), "Register User").Render(r.Context(), w)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal("Error generating hash:", err)
	}

	err = h.repo.CreateUser(r.Context(), auth.FullName, auth.Email, string(passwordHash))
	if err != nil {
		log.Fatal("Cannot create user:", err)
	}
	log.Println("User registered")
}

func (h *AuthHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.LoginForm(model.User{}, map[string]string{}), "Login").Render(r.Context(), w)
}

func (h *AuthHandler) HandleLoginForm(w http.ResponseWriter, r *http.Request) {
	errors := map[string]string{}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := h.repo.GetUser(r.Context(), email)
	if err == repository.ErrNotFound {
		errors["email"] = "User not found"
	} else if err != nil {
		log.Fatal("Error reading user: ", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		errors["password"] = "Invalid email or password"
	} else {
		errors["password"] = "Valid user"
	}

	exp := time.Now().Add(time.Duration(time.Hour * 24))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"iat": time.Now().Unix(),
		"exp": exp.Unix(),
	})

	tokenString, err := token.SignedString(h.hmacTokenSecret)
	if err != nil {
		fmt.Printf("Error creating token: %s\n", err)
		http.Error(w, "Internal error signing in", 500)
		return
	}
	fmt.Println("token: ", tokenString)

	// TODO: Hardening
	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		Expires:  exp,
		MaxAge:   3600 * 24,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/meals", 303)
}
