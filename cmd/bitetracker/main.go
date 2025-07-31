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
	authHandler := handler.NewAuthHandler(&userRepo, config.HmacTokenSecret)
	noAuthMwr := middleware.NewChainNoAuth()
	authMwr := middleware.NewChainWithAuth(config.HmacTokenSecret)

	assetHandler := http.FileServer(http.Dir("internal/view/assets"))

	http.Handle("GET /auth/register", noAuthMwr.Chain(authHandler.RegisterUserForm))
	http.Handle("POST /auth/register", noAuthMwr.Chain(authHandler.HandleRegisterUserForm))
	http.Handle("GET /auth/login", noAuthMwr.Chain(authHandler.LoginForm))
	http.Handle("POST /auth/login", noAuthMwr.Chain(authHandler.HandleLoginForm))
	http.Handle("GET /meals", authMwr.Chain(mealHandler.ListMeals))
	http.Handle("GET /meals/{id}", authMwr.Chain(mealHandler.EditMealForm))
	http.Handle("PUT /meals/{id}", authMwr.Chain(mealHandler.HandleMealForm))
	http.Handle("GET /meals/new", authMwr.Chain(mealHandler.NewMealForm))
	http.Handle("POST /meals/new", authMwr.Chain(mealHandler.HandleMealForm))
	http.Handle("GET /assets/", noAuthMwr.Chain(http.StripPrefix("/assets", assetHandler).ServeHTTP))

	fmt.Println("Server started!")
	http.ListenAndServe(":8000", nil)
}
