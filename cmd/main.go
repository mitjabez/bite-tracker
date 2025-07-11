package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/config"
	db "github.com/mitjabez/bite-tracker/db/init"
	"github.com/mitjabez/bite-tracker/handlers"
	"github.com/mitjabez/bite-tracker/views"
)

func main() {
	config := config.LocalDev()
	dbContext, err := db.Init(config)
	if err != nil {
		log.Fatal("Cannot initialize DB:", err)
	}
	defer dbContext.Pool.Close()

	mealLogHandler := handlers.NewMealLogHandler(dbContext, config.DefaultAppUsername)

	addMealView := views.Base(views.AddMeal(), "Add Meal")
	assetsHandler := http.FileServer(http.Dir("views/assets"))

	http.HandleFunc("/", mealLogHandler.ServeHTTPLogs)
	http.Handle("/add-meal", templ.Handler(addMealView))
	http.Handle("/assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}
