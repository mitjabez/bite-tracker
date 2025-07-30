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

	mealRepo := repository.NewMealRepo(&dbContext)
	userRepo := repository.NewUserRepo(&dbContext)
	mealHandler := handler.NewMealHandler(&mealRepo, config.DefaultAppUserId)
	authHandler := handler.NewAuthHandler(&userRepo)

	assetHandler := http.FileServer(http.Dir("internal/view/assets"))

	http.Handle("GET /auth/register", middleware.Chain(authHandler.RegisterUserForm))
	http.Handle("POST /auth/register", middleware.Chain(authHandler.HandleRegisterUserForm))
	http.Handle("GET /meals", middleware.Chain(mealHandler.ListMeals))
	http.Handle("GET /meals/{id}", middleware.Chain(mealHandler.EditMealForm))
	http.Handle("PUT /meals/{id}", middleware.Chain(mealHandler.HandleMealForm))
	http.Handle("GET /meals/new", middleware.Chain(mealHandler.NewMealForm))
	http.Handle("POST /meals/new", middleware.Chain(mealHandler.HandleMealForm))
	http.Handle("GET /assets/", middleware.Chain(http.StripPrefix("/assets", assetHandler).ServeHTTP))

	fmt.Println("Server started!")
	http.ListenAndServe(":8000", nil)
}
