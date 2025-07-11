package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/internal/config"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/handlers"
	"github.com/mitjabez/bite-tracker/internal/views"
)

func main() {
	config := config.LocalDev()
	dbContext, err := db.Init(config)
	if err != nil {
		log.Fatal("Cannot initialize DB:", err)
	}
	defer dbContext.Pool.Close()

	mealLogHandler := handlers.NewMealHandler(dbContext, config.DefaultAppUsername)

	addMealView := views.Layout(views.MealsNew(), "New Meal")
	assetsHandler := http.FileServer(http.Dir("views/assets"))

	http.HandleFunc("GET /meals", mealLogHandler.ListMeals)
	http.Handle("GET /meals/new", templ.Handler(addMealView))
	http.Handle("GET /assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}
