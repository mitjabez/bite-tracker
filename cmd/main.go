package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/handlers"
	mealservice "github.com/mitjabez/bite-tracker/service"
	"github.com/mitjabez/bite-tracker/views"
)

type Config struct {
	Username         string
	ConnectionString string
}

func main() {
	// doSQL()
	config := Config{
		Username:         "salsajimmy",
		ConnectionString: "postgres://biteapp:superburrito@localhost:5432/bite_tracker?sslmode=disable",
	}
	dbConnection, err := mealservice.New(config.ConnectionString)
	if err != nil {
		log.Fatal("Error initializing DB", err)
	}
	mealLogHandler := handlers.NewMealLogHandler(dbConnection, config.Username)

	// logView := views.Base(views.Log(meals), "Bite Log")
	addMealView := views.Base(views.AddMeal(), "Add Meal")
	assetsHandler := http.FileServer(http.Dir("views/assets"))

	// http.Handle("/", templ.Handler(logView))
	http.HandleFunc("/", mealLogHandler.ServeHTTPLogs)
	http.Handle("/add-meal", templ.Handler(addMealView))
	http.Handle("/assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}

// func doSQL() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	log.Print("Connecting to DB ...")
// 	conn, err := pgx.Connect(ctx, "postgres://biteapp:superburrito@localhost:5432/bite_tracker?sslmode=disable")
// 	if err != nil {
// 		log.Fatal("Cannot open DB:", err)
// 	}
// 	defer conn.Close(ctx)
// 	log.Println("DONE")
//
// 	queries := sqlc.New(conn)
// 	myUUID, err := uuid.Parse("f41ad27a-881d-4f7f-a908-f16a26ce7b78")
// 	if err != nil {
// 		log.Fatal("Error parsing UUID", err)
// 	}
//
// 	log.Print("Querying meals ...")
// 	meals, err := queries.ListMealsByDate(ctx, sqlc.ListMealsByDateParams{
// 		UserID:  myUUID,
// 		ForDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.Now().UTC().Location()),
// 	})
// 	if err != nil {
// 		log.Fatal("Error querying DB:", err)
// 	}
// 	log.Println("Got some meals:", len(meals))
// }
