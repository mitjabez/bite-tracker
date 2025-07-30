package handler

import (
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mitjabez/bite-tracker/internal/model"
	"github.com/mitjabez/bite-tracker/internal/repository"
	"github.com/mitjabez/bite-tracker/internal/view"
)

type Mealhandler struct {
	repo *repository.MealRepo
	// TODO: Move to session
	userId uuid.UUID
}

func NewMealHandler(repo *repository.MealRepo, userId string) Mealhandler {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		log.Fatal("Error parsing uuid", userId, err)
	}
	return Mealhandler{repo: repo, userId: userUUID}
}

func (h Mealhandler) ListMeals(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println("No auth cookie found")
			http.Redirect(w, r, "/auth/login", 302)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	hmacTokenSecret := []byte("1WSB6LaNNLfxi.JbTxrao0s3b4wTpH")
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
		return hmacTokenSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId := claims["sub"]
		exp := claims["exp"]
		log.Println("Got claims:", userId, exp)
	} else {
		log.Fatal("Error obtaining claims")
	}

	date := dateParam(r)
	currentDate := date.Format("2006-01-02")
	prevDate := date.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	mealsView, err := h.repo.ListMeals(r.Context(), h.userId, date)
	if err != nil {
		log.Fatal("Error querying users")
	}
	view.Layout(view.ListMeals(prevDate, nextDate, currentDate, mealsView), "Meal Log").Render(r.Context(), w)
}

func (h Mealhandler) NewMealForm(w http.ResponseWriter, r *http.Request) {
	date := dateParam(r)
	now := time.Now()
	mealView := model.Meal{
		TimeOfMeal:  time.Date(date.Year(), date.Month(), date.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location()),
		HungerLevel: 4,
	}

	top3Meals, err := h.repo.Top3Meals(r.Context(), h.userId)
	if err != nil {
		log.Fatal("Error retrieving top meals for user: ", err)
	}

	view.Layout(view.NewMealForm(mealView, map[string]string{}, model.Symptoms, top3Meals), "New Meal").Render(r.Context(), w)
}

func (h Mealhandler) EditMealForm(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.PathValue("id")
	userUUID, err := uuid.Parse(userIdParam)
	if err != nil {
		log.Fatal("Invalid uuid ", userIdParam)
	}

	mealView, err := h.repo.GetMeal(r.Context(), userUUID)
	if err != nil {
		log.Fatal("Error reading meal: ", err)
	}

	top3Meals, err := h.repo.Top3Meals(r.Context(), userUUID)
	if err != nil {
		log.Fatal("Error retrieving top meals for user: ", err)
	}
	view.Layout(view.EditMealForm(mealView, map[string]string{}, model.Symptoms, top3Meals), "Edit Meal").Render(r.Context(), w)
}

func (h Mealhandler) HandleMealForm(w http.ResponseWriter, r *http.Request) {
	mealIdParam := r.PathValue("id")
	var mealUUID uuid.UUID
	var err error
	isNewMeal := r.Method == "POST"
	if !isNewMeal {
		mealUUID, err = uuid.Parse(mealIdParam)
		if err != nil {
			log.Fatal("Invalid meal uuid ", mealIdParam)
		}
	}

	dateParam := r.FormValue("date")
	timeParam := r.FormValue("time")
	mealParam := strings.Trim(r.FormValue("meal"), " ")
	hungerParam := r.FormValue("hunger")
	symptoms := r.PostForm["symptoms"]

	errors := map[string]string{}

	if len(mealParam) == 0 {
		errors["meal"] = "Meal is required"
	}

	_, err = time.Parse("2006-01-02", dateParam)
	if err != nil {
		errors["date"] = "Invalid date"
	}

	_, err = time.Parse("15:04", timeParam)
	if err != nil {
		errors["time"] = "Invalid time"
	}
	dateAndTime, err := time.Parse("2006-01-02 15:04", dateParam+" "+timeParam)
	if err != nil {
		errors["time"] = "Invalid date or time"
		errors["date"] = "Invalid date or time"
	}

	if len(mealParam) == 0 {
		errors["meal"] = "Meal is required"
	}

	hungerLevel, err := strconv.Atoi(hungerParam)
	if err != nil || hungerLevel < 1 || hungerLevel > 5 {
		errors["hunger"] = "Invalid hunger level"
	}

	for _, s := range symptoms {
		if !slices.Contains(model.Symptoms, s) {
			errors["symptoms"] = "Symptom " + s + " doesn't exist"
			break
		}
	}

	mealView := model.Meal{
		Id:          mealIdParam,
		MealType:    model.ResolveMealType(dateAndTime),
		TimeOfMeal:  dateAndTime,
		Description: mealParam,
		HungerLevel: int32(hungerLevel),
		Symptoms:    symptoms,
	}

	if len(errors) > 0 {
		top3Meals, err := h.repo.Top3Meals(r.Context(), h.userId)
		if err != nil {
			log.Fatal("Error obtaining top 3 meals")
		}
		if isNewMeal {
			view.Layout(view.NewMealForm(mealView, errors, model.Symptoms, top3Meals), "New Meal").Render(r.Context(), w)
		} else {
			view.Layout(view.EditMealForm(mealView, errors, model.Symptoms, top3Meals), "Edit Meal").Render(r.Context(), w)
		}
		return
	}

	if isNewMeal {
		h.repo.CreateMeal(r.Context(), h.userId, mealView)
	} else {
		h.repo.UpdateMeal(r.Context(), h.userId, mealUUID, mealView)
	}
	if err != nil {
		log.Fatal("Cannot create or update meal: ", err)
	}

	http.Redirect(w, r, "/meals", 303)
}

func dateParam(r *http.Request) time.Time {
	dateQuery := r.FormValue("date")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if dateQuery != "" {
		parsedDate, err := time.Parse("2006-01-02", dateQuery)
		if err != nil {
			log.Println("WARNING: Error parsing date", dateQuery)
		} else {
			date = parsedDate
		}
	}
	return date
}
