package handler

import (
	"net/http"
	"net/mail"
	"regexp"
	"strconv"

	"github.com/mitjabez/bite-tracker/internal/auth"
	"github.com/mitjabez/bite-tracker/internal/httpx"
	"github.com/mitjabez/bite-tracker/internal/model"
	"github.com/mitjabez/bite-tracker/internal/repository"
	"github.com/mitjabez/bite-tracker/internal/view"
	"golang.org/x/crypto/bcrypt"
)

var (
	lowercaseRegex = regexp.MustCompile(".*[a-z].*")
	uppercaseRegex = regexp.MustCompile(".*[A-Z].*")
	numbersRegex   = regexp.MustCompile(".*[0-9].*")
)

const minPasswordLen = 10
const maxPasswordLen = 100

type AuthHandler struct {
	repo *repository.UserRepo
	auth *auth.Auth
}

func NewAuthHandler(repo *repository.UserRepo, auth *auth.Auth) *AuthHandler {
	return &AuthHandler{
		repo: repo,
		auth: auth,
	}
}

func (h *AuthHandler) RegisterUserForm(w http.ResponseWriter, r *http.Request) {
	view.NotLoggedInLayout(view.RegisterUserForm(model.User{}, "", "", map[string]string{}), "Register User").Render(r.Context(), w)
}

func (h *AuthHandler) UserProfileForm(w http.ResponseWriter, r *http.Request, dbUser model.User) {
	// Get user from DB just in case additional fields got added not yet in the token
	dbUser, err := h.repo.GetUser(r.Context(), dbUser.Id)
	if err != nil {
		httpx.InternalError(w, "Failed obtaining user", err)
		return
	}
	view.LoggedInLayout(view.UserProfileForm(dbUser, "", "", map[string]string{}), "User Profile", dbUser).Render(r.Context(), w)
}

func (h *AuthHandler) HandleUserProfileForm(w http.ResponseWriter, r *http.Request, dbUser model.User) {
	// If we would go directly auth would not be checked
	h.handleUserForm(w, r, model.User{})
}

func (h *AuthHandler) HandleRegisterUserForm(w http.ResponseWriter, r *http.Request) {
	h.handleUserForm(w, r, model.User{})
}

func (h *AuthHandler) handleUserForm(w http.ResponseWriter, r *http.Request, user model.User) {
	errors := map[string]string{}
	// TODO: Check also for PUT
	isNewUser := r.Method == "POST"

	userForm := model.User{
		FullName: r.FormValue("full-name"),
		Email:    r.FormValue("email"),
	}
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	if len(userForm.FullName) < 5 {
		errors["full-name"] = "Full name must be at least 5 characters long."
	}

	_, err := mail.ParseAddress(userForm.Email)
	if err != nil {
		errors["email"] = "Invalid email address"
	}

	if !verifyPasswordComplexity(password) {
		errors["password"] = "Invalid password. It must be between " + strconv.Itoa(minPasswordLen) + " and " +
			strconv.Itoa(minPasswordLen) + " characters long. " +
			"Must include at least one lowercase letter, one uppercase letter and one number."
	}

	// Show only one password message at a time
	if errors["password"] == "" && password != confirmPassword {
		errors["confirmPassword"] = "Passwords do not match"
	}

	if isNewUser && len(errors) == 0 {
		userExists, err := h.repo.UserExists(r.Context(), userForm.Email)
		if err != nil {
			httpx.InternalError(w, "Failed checking if user exists", err)
			return
		}
		if userExists {
			errors["email"] = "Username unavailable — please choose another"
		}
	}

	if len(errors) > 0 {
		if isNewUser {
			view.NotLoggedInLayout(view.RegisterUserForm(userForm, password, confirmPassword, errors), "Register User").Render(r.Context(), w)
		} else {
			view.LoggedInLayout(view.RegisterUserForm(userForm, password, confirmPassword, errors), "Register User", user).Render(r.Context(), w)
		}
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		httpx.InternalError(w, "Failed generating hash", err)
		return
	}

	if isNewUser {
		_, err := h.repo.CreateUser(r.Context(), userForm.FullName, userForm.Email, string(passwordHash))
		if err != nil {
			httpx.InternalError(w, "Failed creating user", err)
			return
		}
		view.NotLoggedInLayout(view.ProfileUpdated("Registration successful!"), "Profile Created").Render(r.Context(), w)
	} else {
		err = h.repo.UpdateUser(r.Context(), user.Id, userForm.FullName, userForm.Email, string(passwordHash))
		if err != nil {
			httpx.InternalError(w, "Failed updating user", err)
			return
		}
		h.auth.InvalidateCookieToken(w)
		view.NotLoggedInLayout(view.ProfileUpdated("Your profile has been successfully updated!"), "Profile Updated").Render(r.Context(), w)
	}
}

func (h *AuthHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	view.NotLoggedInLayout(view.LoginForm(model.User{}, map[string]string{}), "Login").Render(r.Context(), w)
}

func (h *AuthHandler) HandleLoginForm(w http.ResponseWriter, r *http.Request) {
	errors := map[string]string{}
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(email) < 5 || len(email) > 100 {
		errors["email"] = "Email must be between 5 and 100 characters long"
	}
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		errors["password"] = "Password must be between " + strconv.Itoa(minPasswordLen) + " and " + strconv.Itoa(maxPasswordLen) + " characters long"
	}

	if len(errors) > 0 {
		handleInvalidLogin(errors, email, w, r)
		return
	}

	user, err := h.repo.GetUserByEmail(r.Context(), email)
	if err == repository.ErrNotFound {
		errors["email"] = "Invalid email or password"
		handleInvalidLogin(errors, email, w, r)
		return
	} else if err != nil {
		httpx.InternalError(w, "Failed reading user ", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		errors["email"] = "Invalid email or password"
		handleInvalidLogin(errors, email, w, r)
		return
	}

	err = h.auth.SetCookieToken(w, user)
	if err != nil {
		httpx.InternalError(w, "Failed setting cookie token ", err)
		return
	}
	http.Redirect(w, r, "/meals", 302)
}

func verifyPasswordComplexity(password string) bool {
	return len(password) >= minPasswordLen &&
		len(password) <= maxPasswordLen &&
		lowercaseRegex.MatchString(password) &&
		uppercaseRegex.MatchString(password) &&
		numbersRegex.MatchString(password)
}

func handleInvalidLogin(errors map[string]string, email string, w http.ResponseWriter, r *http.Request) {
	view.NotLoggedInLayout(view.LoginForm(model.User{Email: email}, errors), "Login").Render(r.Context(), w)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	h.auth.InvalidateCookieToken(w)
	redirectToLogin(w, r)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", 302)
}
