package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mitjabez/bite-tracker/db/sqlc"
	"github.com/mitjabez/bite-tracker/models"
	"github.com/mitjabez/bite-tracker/views"
)

func main() {
	meals := []models.Meal{
		{
			Id:          "1",
			Type:        "Breakfast",
			Time:        time.Now(),
			Description: "Classic yoghurt with oats",
			HungerLevel: 4,
		},
		{
			Id:          "2",
			Type:        "Lunch",
			Time:        time.Now().Add(-4 * time.Hour),
			Description: "Grilled chicken with mixed greens and quinoa",
			HungerLevel: 5,
			Symptoms:    []models.MealSymptom{models.Bloating},
		},
		{
			Id:          "3",
			Type:        "Snack",
			Time:        time.Now().Add(-2 * time.Hour),
			Description: "Handful of almonds and an apple",
			HungerLevel: 3,
		},
		{
			Id:          "4",
			Type:        "Dinner",
			Time:        time.Now().Add(-8 * time.Hour),
			Description: "Salmon with roasted sweet potatoes and broccoli",
			HungerLevel: 5,
		},
	}

	doSQL()

	logView := views.Base(views.Log(meals), "Bite Log")
	addMealView := views.Base(views.AddMeal(), "Add Meal")
	assetsHandler := http.FileServer(http.Dir("views/assets"))

	http.Handle("/", templ.Handler(logView))
	http.Handle("/add-meal", templ.Handler(addMealView))
	http.Handle("/assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}

func doSQL() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, "postgres://biteapp:superburrito@localhost:5432/bite_tracker?sslmode=disable")
	if err != nil {
		log.Fatal("Cannot open DB:", err)
	}
	defer conn.Close(ctx)

	queries := sqlc.New(conn)
	meals, err := queries.ListMeals(ctx, pgtype.Timestamptz{})
	if err != nil {
		log.Fatal("Error querying DB:", err)
	}
	log.Println("Got some meals:", len(meals))
}
