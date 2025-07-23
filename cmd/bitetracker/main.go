package main

import (
	"log"
	"net/http"

	"github.com/mitjabez/bite-tracker/internal/config"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/handlers"
)

func main() {
	config := config.LocalDev()
	dbContext, err := db.Init(config)
	if err != nil {
		log.Fatal("Cannot initialize DB:", err)
	}
	defer dbContext.Pool.Close()

	mealLogHandler := handlers.NewMealHandler(dbContext, config.DefaultAppUserId)

	assetsHandler := http.FileServer(http.Dir("views/assets"))

	http.HandleFunc("GET /meals/{id}", mealLogHandler.EditMeal)
	http.HandleFunc("GET /meals", mealLogHandler.ListMeals)
	http.HandleFunc("GET /meals/new", mealLogHandler.NewMeal)
	http.HandleFunc("PUT /meals/{id}/edit", mealLogHandler.HandleMealForm)
	http.HandleFunc("POST /meals/new", mealLogHandler.HandleMealForm)
	http.Handle("GET /assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}
