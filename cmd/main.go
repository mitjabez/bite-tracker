package main

import (
	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/views"
	"net/http"
)

func main() {
	logView := views.Base(views.Log(), "Bite Log")
	addMealView := views.Base(views.AddMeal(), "Add meal")
	assetsHandler := http.FileServer(http.Dir("views/assets"))

	http.Handle("/", templ.Handler(logView))
	http.Handle("/add-meal", templ.Handler(addMealView))
	http.Handle("/assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}
