package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
	"github.com/mitjabez/bite-tracker/internal/models"
	"github.com/mitjabez/bite-tracker/internal/views"
)

type MealLogHandler struct {
	queries  *sqlc.Queries
	username string
}

func NewMealLogHandler(dbContext db.DBContext, username string) MealLogHandler {
	return MealLogHandler{queries: dbContext.Queries, username: username}
}

func (mealLogHandler MealLogHandler) ServeHTTPLogs(w http.ResponseWriter, r *http.Request) {
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

	currentDate := date.Format("2006-01-02")
	prevDate := date.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	params := sqlc.ListMealsByUsernameAndDateParams{
		Username: mealLogHandler.username,
		ForDate:  date,
	}
	meals, err := mealLogHandler.queries.ListMealsByUsernameAndDate(ctx, params)
	if err != nil {
		log.Fatal("Error retrieving meals", err)
	}

	mealsView := []models.MealView{}
	for _, m := range meals {
		usedSymptoms, unusedSymptoms := splitSymptoms(m.Symptoms)
		mealsView = append(mealsView, models.MealView{
			MealType:       m.MealType,
			TimeOfMeal:     m.TimeOfMeal.Format("15:04"),
			Description:    m.Description,
			HungerLevel:    m.HungerLevel,
			UsedSymptoms:   usedSymptoms,
			UnusedSymptoms: unusedSymptoms,
		})
	}

	views.Base(views.Log(prevDate, nextDate, currentDate, mealsView), "Meal Log").Render(r.Context(), w)
}

func splitSymptoms(symptoms []string) (usedSymptoms, unusedSymptoms []models.MealSymptom) {
	usedSymptomNames := map[string]bool{}
	for _, name := range symptoms {
		usedSymptomNames[name] = true
	}

	for name, symptom := range models.AllSymptoms {
		if usedSymptomNames[name] {
			usedSymptoms = append(usedSymptoms, symptom)
		} else {
			unusedSymptoms = append(unusedSymptoms, symptom)
		}
	}
	return
}
