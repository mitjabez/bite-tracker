package handler

import (
	"log"
	"net/http"
	"net/mail"

	"github.com/mitjabez/bite-tracker/internal/model"
	"github.com/mitjabez/bite-tracker/internal/repository"
	"github.com/mitjabez/bite-tracker/internal/view"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo *repository.UserRepo
}

func NewAuthHandler(repo *repository.UserRepo) AuthHandler {
	return AuthHandler{
		repo: repo,
	}
}

func (h *AuthHandler) RegisterUserForm(w http.ResponseWriter, r *http.Request) {
	view.Layout(view.RegisterUserForm(model.User{}, map[string]string{}), "Register User").Render(r.Context(), w)
}

func (h *AuthHandler) HandleRegisterUserForm(w http.ResponseWriter, r *http.Request) {
	errors := map[string]string{}

	auth := model.User{
		FullName: r.FormValue("full-name"),
		EMail:    r.FormValue("email"),
	}
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	if len(auth.FullName) < 5 {
		errors["full-name"] = "Full name must be at least 5 characters long."
	}

	_, err := mail.ParseAddress(auth.EMail)
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
		userExists, err := h.repo.UserExists(r.Context(), auth.EMail)
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

	err = h.repo.CreateUser(r.Context(), auth.FullName, auth.EMail, string(passwordHash))
	if err != nil {
		log.Fatal("Cannot create user:", err)
	}
	log.Println("User registered")

}
