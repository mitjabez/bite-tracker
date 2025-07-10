package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mitjabez/bite-tracker/db/sqlc"
	"github.com/mitjabez/bite-tracker/models"
	"github.com/mitjabez/bite-tracker/views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ServeHTTPLogs(writer http.ResponseWriter, request *http.Request) {
	meals := doSQL()
	views.Base(views.Log(meals), "Meal Log").Render(request.Context(), writer)
}

func doSQL() []models.Meal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Print("Connecting to DB ...")
	conn, err := pgx.Connect(ctx, "postgres://biteapp:superburrito@localhost:5432/bite_tracker?sslmode=disable")
	if err != nil {
		log.Fatal("Cannot open DB:", err)
	}
	defer conn.Close(ctx)
	log.Println("DONE")

	queries := sqlc.New(conn)
	myUUID, err := uuid.Parse("f41ad27a-881d-4f7f-a908-f16a26ce7b78")
	if err != nil {
		log.Fatal("Error parsing UUID", err)
	}

	log.Print("Querying meals ...")
	meals, err := queries.ListMealsByDate(ctx, sqlc.ListMealsByDateParams{
		UserID:  myUUID,
		ForDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.Now().UTC().Location()),
	})
	if err != nil {
		log.Fatal("Error querying DB:", err)
	}
	log.Println("Got some meals:", len(meals))

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
	return httpMeals
}
