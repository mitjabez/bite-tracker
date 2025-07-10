package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/mitjabez/bite-tracker/models"
	mealservice "github.com/mitjabez/bite-tracker/service"
	"github.com/mitjabez/bite-tracker/views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TODO: Conolidate log/meallog
type MealLogHandler struct {
	dbConnection mealservice.DBConnection
}

func NewMealLogHandler() MealLogHandler {
	dbconn, err := mealservice.New()
	if err != nil {
		log.Fatal("Error initializing DB", err)
	}
	return MealLogHandler{dbConnection: dbconn}

}

func (mealLogHandler MealLogHandler) ServeHTTPLogs(writer http.ResponseWriter, request *http.Request) {
	meals, err := mealLogHandler.dbConnection.GetMeals("f41ad27a-881d-4f7f-a908-f16a26ce7b78", time.Date(2025, 3, 1, 0, 0, 0, 0, time.Now().UTC().Location()))
	if err != nil {
		log.Fatal("Error getting meals")
	}
	httpMeals := []models.Meal{}
	for _, m := range meals {
		symptoms := []models.MealSymptom{}
		for _, s := range m.Symptoms {
			symptoms = append(symptoms, models.MealSymptom(s))
		}
		meal := models.Meal{
			Id:          m.ID.String(),
			Type:        cases.Title(language.English, cases.Compact).String(m.MealType),
			Time:        m.TimeOfMeal,
			Description: m.Description,
			HungerLevel: int64(m.HungerLevel),
			Symptoms:    symptoms,
		}

		httpMeals = append(httpMeals, meal)
	}

	views.Base(views.Log(httpMeals), "Meal Log").Render(request.Context(), writer)
}
