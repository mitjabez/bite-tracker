package main

import (
	"log"
	"net/http"

	"github.com/mitjabez/bite-tracker/internal/auth"
	"github.com/mitjabez/bite-tracker/internal/config"
	"github.com/mitjabez/bite-tracker/internal/db"
	"github.com/mitjabez/bite-tracker/internal/handler"
	"github.com/mitjabez/bite-tracker/internal/middleware"
	"github.com/mitjabez/bite-tracker/internal/repository"
)

func main() {
	config := config.LocalDev()
	dbContext, err := db.Init(config)
	if err != nil {
		log.Fatal("Failed initializing DB: ", err)
	}
	defer dbContext.Pool.Close()

	mealRepo := repository.NewMealRepo(&dbContext)
	userRepo := repository.NewUserRepo(&dbContext)
	auth := auth.NewAuth(config.HmacTokenSecret, config.TokenAge)
	mealHandler := handler.NewMealHandler(mealRepo, auth)
	authHandler := handler.NewAuthHandler(userRepo, auth)
	mwr := middleware.New(auth)

	assetHandler := http.FileServer(http.Dir("internal/view/assets"))

	http.Handle("GET /", mwr.Chain(handler.Home))
	http.Handle("GET /auth/register", mwr.Chain(authHandler.RegisterUserForm))
	http.Handle("POST /auth/register", mwr.Chain(authHandler.HandleRegisterUserForm))
	http.Handle("GET /auth/login", mwr.Chain(authHandler.LoginForm))
	http.Handle("POST /auth/login", mwr.Chain(authHandler.HandleLoginForm))
	http.Handle("GET /auth/logout", mwr.Chain(authHandler.HandleLogout))
	http.Handle("GET /auth/profile", mwr.AuthChain(authHandler.UserProfileForm))
	http.Handle("PUT /auth/profile", mwr.AuthChain(authHandler.HandleUserProfileForm))
	http.Handle("GET /meals", mwr.AuthChain(mealHandler.ListMeals))
	http.Handle("GET /meals/{id}", mwr.AuthChain(mealHandler.EditMealForm))
	http.Handle("PUT /meals/{id}", mwr.AuthChain(mealHandler.HandleMealForm))
	http.Handle("DELETE /meals/{id}", mwr.AuthChain(mealHandler.HandleDelete))
	http.Handle("GET /meals/new", mwr.AuthChain(mealHandler.NewMealForm))
	http.Handle("POST /meals/new", mwr.AuthChain(mealHandler.HandleMealForm))
	http.Handle("GET /assets/", mwr.Chain(http.StripPrefix("/assets", assetHandler).ServeHTTP))

	log.Printf("Bite Tracker started on %s\n", config.ListenAddr)
	err = http.ListenAndServe(config.ListenAddr, nil)
	if err != nil {
		log.Fatalf("Failed starting server: %v\n", err)
	}
}
