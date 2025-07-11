package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/mitjabez/bite-tracker/db/init"
	"github.com/mitjabez/bite-tracker/db/sqlc"
	"github.com/mitjabez/bite-tracker/views"
)

type MealLogHandler struct {
	queries  *sqlc.Queries
	username string
}

func NewMealLogHandler(dbContext db.DBContext, username string) MealLogHandler {
	return MealLogHandler{queries: dbContext.Queries, username: username}
}

func (mealLogHandler MealLogHandler) ServeHTTPLogs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	params := sqlc.ListMealsByUsernameAndDateParams{
		Username: mealLogHandler.username,
		ForDate:  time.Date(2025, 3, 1, 0, 0, 0, 0, time.Now().UTC().Location()),
	}
	meals, err := mealLogHandler.queries.ListMealsByUsernameAndDate(ctx, params)
	if err != nil {
		log.Fatal("Error retrieving meals", err)
	}

	views.Base(views.Log(meals), "Meal Log").Render(r.Context(), w)
}
