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
	btConfig, err := config.Init()
	if err != nil {
		log.Fatal("Cannot load app config: ", err)
	}

	dbContext, err := db.Init(btConfig)
	if err != nil {
		log.Fatal("Failed initializing DB: ", err)
	}
	defer dbContext.Pool.Close()
	err = db.RunMigration(btConfig)
	if err != nil {
		log.Printf("Error running DB migration. App may not run correctly. Error: %v\n", err)
	}

	mealRepo := repository.NewMealRepo(&dbContext)
	userRepo := repository.NewUserRepo(&dbContext)
	auth := auth.NewAuth(btConfig.HmacTokenSecret, btConfig.TokenAge)
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

	log.Printf("Bite Tracker started on %s\n", btConfig.ListenAddr)
	err = http.ListenAndServe(btConfig.ListenAddr, nil)
	if err != nil {
		log.Fatalf("Failed starting server: %v\n", err)
	}
}
