package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mitjabez/bite-tracker/internal/config"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/handler"
	"github.com/mitjabez/bite-tracker/internal/middleware"
	"github.com/mitjabez/bite-tracker/internal/repository"
)

func main() {
	config := config.LocalDev()
	dbContext, err := db.Init(config)
	if err != nil {
		log.Fatal("Cannot initialize DB:", err)
	}
	defer dbContext.Pool.Close()

	repository := repository.MealRepo{
		DBContext: dbContext,
	}
	mealLogHandler := handler.NewMealHandler(&repository, config.DefaultAppUserId)

	assetsHandler := http.FileServer(http.Dir("internal/view/assets"))

	http.Handle("GET /meals", middleware.Logger(http.HandlerFunc(mealLogHandler.ListMeals)))
	http.Handle("GET /meals/{id}", middleware.Logger(http.HandlerFunc(mealLogHandler.EditMeal)))
	http.Handle("PUT /meals/{id}", middleware.Logger(http.HandlerFunc(mealLogHandler.HandleMealForm)))
	http.Handle("GET /meals/new", middleware.Logger(http.HandlerFunc(mealLogHandler.NewMeal)))
	http.Handle("POST /meals/new", middleware.Logger(http.HandlerFunc(mealLogHandler.HandleMealForm)))
	http.Handle("GET /assets/", middleware.Logger(http.StripPrefix("/assets", assetsHandler)))

	fmt.Println("Server started!")
	http.ListenAndServe(":8000", nil)
}
