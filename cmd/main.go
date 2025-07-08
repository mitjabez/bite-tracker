package main

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
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

	logView := views.Base(views.Log(meals), "Bite Log")
	addMealView := views.Base(views.AddMeal(), "Add Meal")
	assetsHandler := http.FileServer(http.Dir("views/assets"))

	http.Handle("/", templ.Handler(logView))
	http.Handle("/add-meal", templ.Handler(addMealView))
	http.Handle("/assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}
