package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
	"github.com/mitjabez/bite-tracker/internal/models"
	"github.com/mitjabez/bite-tracker/internal/views"
)

type Mealhandler struct {
	queries *sqlc.Queries
	// TODO: Move to session
	userId uuid.UUID
}

func NewMealHandler(dbContext db.DBContext, userId string) Mealhandler {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		log.Fatal("Error parsing uuid", userId, err)
	}
	return Mealhandler{queries: dbContext.Queries, userId: userUUID}
}

func (h Mealhandler) ListMeals(w http.ResponseWriter, r *http.Request) {
	date := dateParam(r)
	currentDate := date.Format("2006-01-02")
	prevDate := date.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	params := sqlc.ListMealsByUsernameAndDateParams{
		UserID:  h.userId,
		ForDate: date,
	}
	meals, err := h.queries.ListMealsByUsernameAndDate(ctx, params)
	if err != nil {
		log.Fatal("Error retrieving meals: ", err)
	}

	mealsView := []models.MealView{}
	for _, m := range meals {
		unusedSymptoms := getUnusedSymptoms(m.Symptoms)
		mealsView = append(mealsView, models.MealView{
			MealType:       m.MealTypeID,
			DateOfMeal:     m.TimeOfMeal.Format("2006-01-02"),
			TimeOfMeal:     m.TimeOfMeal.Format("15:04"),
			Description:    m.Description,
			HungerLevel:    m.HungerLevel,
			UsedSymptoms:   m.Symptoms,
			UnusedSymptoms: unusedSymptoms,
		})
	}

	views.Layout(views.Meals(prevDate, nextDate, currentDate, mealsView), "Meal Log").Render(r.Context(), w)
}

func (h Mealhandler) NewMeal(w http.ResponseWriter, r *http.Request) {
	mealView := models.MealView{
		DateOfMeal:  dateParam(r).Format("2006-01-02"),
		TimeOfMeal:  time.Now().Format("15:04"),
		HungerLevel: 4,
	}
	views.Layout(views.MealsNew(mealView, map[string]string{}, models.Symptoms), "New Meal").Render(r.Context(), w)
}

func (h Mealhandler) CreateMeal(w http.ResponseWriter, r *http.Request) {
	dateParam := r.FormValue("date")
	timeParam := r.FormValue("time")
	mealParam := strings.Trim(r.FormValue("meal"), " ")
	hungerParam := r.FormValue("hunger")
	symptoms := r.PostForm["symptoms"]

	errors := map[string]string{}

	if len(mealParam) == 0 {
		errors["meal"] = "Meal is required"
	}

	_, err := time.Parse("2006-01-02", dateParam)
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

	mealsView := models.MealView{
		DateOfMeal:  dateParam,
		TimeOfMeal:  timeParam,
		Description: mealParam,
		HungerLevel: int32(hungerLevel),
	}

	if len(errors) > 0 {
		views.Layout(views.MealsNew(mealsView, errors, models.Symptoms), "New Meal").Render(r.Context(), w)
		return
	}

	mealType := "Dinner"
	hour := dateAndTime.Hour()
	switch {
	case hour < 9:
		mealType = "Breakfast"
	case hour < 11:
		mealType = "Brunch"
	case hour < 15:
		mealType = "Lunch"
	default:
		mealType = "Dinner"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	_, err = h.queries.CreateMeal(ctx, sqlc.CreateMealParams{
		UserID:      h.userId,
		MealTypeID:  mealType,
		TimeOfMeal:  dateAndTime,
		Description: mealParam,
		HungerLevel: int32(hungerLevel),
		Symptoms:    symptoms,
	})
	if err != nil {
		log.Fatal("Cannot create meal:", err)
	}

	h.ListMeals(w, r)
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

// TODO: Optimize
func getUnusedSymptoms(activeSymptoms []string) []string {
	usedSymptomNames := map[string]bool{}
	for _, s := range activeSymptoms {
		usedSymptomNames[s] = true
	}

	unusedSymptoms := []string{}
	for _, s := range models.Symptoms {
		if !usedSymptomNames[s] {
			unusedSymptoms = append(unusedSymptoms, s)
		}
	}
	return unusedSymptoms
}
