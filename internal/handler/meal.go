package handler

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitjabez/bite-tracker/internal/auth"
	"github.com/mitjabez/bite-tracker/internal/httpx"
	"github.com/mitjabez/bite-tracker/internal/model"
	"github.com/mitjabez/bite-tracker/internal/repository"
	"github.com/mitjabez/bite-tracker/internal/view"
)

type Mealhandler struct {
	repo *repository.MealRepo
	auth *auth.Auth
}

func NewMealHandler(repo *repository.MealRepo, auth *auth.Auth) *Mealhandler {
	return &Mealhandler{repo: repo, auth: auth}
}

func (h Mealhandler) ListMeals(w http.ResponseWriter, r *http.Request, user model.User) {
	date := dateParam(r)
	currentDate := date.Format("2006-01-02")
	prevDate := date.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	mealsView, err := h.repo.ListMeals(r.Context(), user.Id, date)
	if err != nil {
		httpx.InternalError(w, "Failed to list meals", err)
		return
	}
	view.LoggedInLayout(view.ListMeals(prevDate, nextDate, currentDate, mealsView), "Meal Log", user).Render(r.Context(), w)
}

func (h Mealhandler) NewMealForm(w http.ResponseWriter, r *http.Request, user model.User) {
	date := dateParam(r)
	now := time.Now()

	mealView := model.Meal{
		TimeOfMeal:  time.Date(date.Year(), date.Month(), date.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location()),
		HungerLevel: 4,
	}

	top3Meals, err := h.repo.Top3Meals(r.Context(), user.Id)
	if err != nil {
		httpx.InternalError(w, "Failed to retrieve top meals", err)
		return
	}

	view.LoggedInLayout(view.NewMealForm(mealView, map[string]string{}, model.Symptoms, top3Meals), "New Meal", user).Render(r.Context(), w)
}

func (h Mealhandler) EditMealForm(w http.ResponseWriter, r *http.Request, user model.User) {
	mealIdParam := r.PathValue("id")
	mealUUID, err := uuid.Parse(mealIdParam)
	if err != nil {
		httpx.InternalError(w, "Invalid uuid", err)
		return
	}

	mealView, err := h.repo.GetMeal(r.Context(), mealUUID)
	if err != nil {
		httpx.InternalError(w, "Failed reading meal", err)
		return
	}

	top3Meals, err := h.repo.Top3Meals(r.Context(), mealUUID)
	if err != nil {
		httpx.InternalError(w, "Failed retrieving top meals", err)
		return
	}
	view.LoggedInLayout(view.EditMealForm(mealView, map[string]string{}, model.Symptoms, top3Meals), "Edit Meal", user).Render(r.Context(), w)
}

func (h Mealhandler) HandleMealForm(w http.ResponseWriter, r *http.Request, user model.User) {
	var mealUUID uuid.UUID
	var err error
	mealIdParam := r.PathValue("id")

	isNewMeal := r.Method == "POST"
	if !isNewMeal {
		mealUUID, err = uuid.Parse(mealIdParam)
		if err != nil {
			httpx.InternalError(w, "Invalid meal uuid", err)
			return
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
		top3Meals, err := h.repo.Top3Meals(r.Context(), user.Id)
		if err != nil {
			httpx.InternalError(w, "Failed obtaining top 3 meals", err)
			return
		}
		if isNewMeal {
			view.LoggedInLayout(view.NewMealForm(mealView, errors, model.Symptoms, top3Meals), "New Meal", user).Render(r.Context(), w)
		} else {
			view.LoggedInLayout(view.EditMealForm(mealView, errors, model.Symptoms, top3Meals), "Edit Meal", user).Render(r.Context(), w)
		}
		return
	}

	if isNewMeal {
		h.repo.CreateMeal(r.Context(), user.Id, mealView)
	} else {
		h.repo.UpdateMeal(r.Context(), user.Id, mealUUID, mealView)
	}
	if err != nil {
		httpx.InternalError(w, "Failed creating or update meal: ", err)
		return
	}

	http.Redirect(w, r, "/meals", 303)
}

func (h Mealhandler) HandleDelete(w http.ResponseWriter, r *http.Request, user model.User) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.BadRequest(w, "Invalid uuid", err)
		return
	}

	httpx.InternalError(w, "Cannot delete meal", err)
	err = h.repo.DeleteMeal(r.Context(), id)
	if err != nil {
		httpx.InternalError(w, "Cannot delete meal", err)
		return
	}
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
