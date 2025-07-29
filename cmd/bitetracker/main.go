package main

import (
	"fmt"
	"log"
	"net/http"

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

	mealRepo := repository.MealRepo{DBContext: dbContext}
	userRepo := repository.UserRepo{DBContext: dbContext}
	mealLogHandler := handler.NewMealHandler(&mealRepo, config.DefaultAppUserId)
	authHandler := handler.NewAuthHandler(&userRepo)

	assetsHandler := http.FileServer(http.Dir("internal/view/assets"))

	http.Handle("GET /auth/register", middleware.All(http.HandlerFunc(authHandler.GetRegisterForm)))
	http.Handle("POST /auth/register", middleware.All(http.HandlerFunc(authHandler.PostRegisterForm)))
	http.Handle("GET /meals", middleware.All(http.HandlerFunc(mealLogHandler.ListMeals)))
	http.Handle("GET /meals/{id}", middleware.All(http.HandlerFunc(mealLogHandler.EditMeal)))
	http.Handle("PUT /meals/{id}", middleware.All(http.HandlerFunc(mealLogHandler.HandleMealForm)))
	http.Handle("GET /meals/new", middleware.All(http.HandlerFunc(mealLogHandler.NewMeal)))
	http.Handle("POST /meals/new", middleware.All(http.HandlerFunc(mealLogHandler.HandleMealForm)))
	http.Handle("GET /assets/", middleware.All(http.StripPrefix("/assets", assetsHandler)))

	fmt.Println("Server started!")
	http.ListenAndServe(":8000", nil)
}
