package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/handlers"
	"github.com/mitjabez/bite-tracker/views"
)

func main() {
	// doSQL()

	// logView := views.Base(views.Log(meals), "Bite Log")
	addMealView := views.Base(views.AddMeal(), "Add Meal")
	assetsHandler := http.FileServer(http.Dir("views/assets"))
	mealLogHandler := handlers.NewMealLogHandler()

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
