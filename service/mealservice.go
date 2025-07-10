package mealservice

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mitjabez/bite-tracker/db/sqlc"
	"github.com/mitjabez/bite-tracker/models"
)

type MealService interface {
	GetMeals(userId string, date time.Time) []models.Meal
}

type DBConnection struct {
	ctx     context.Context
	queries *sqlc.Queries
}

func New() (DBConnection, error) {
	// TODO: Handle timeout
	ctx := context.Background() //context.WithTimeout(context.Background(), 5*time.Second)

	log.Print("Connecting to DB ...")
	conn, err := pgx.Connect(ctx, "postgres://biteapp:superburrito@localhost:5432/bite_tracker?sslmode=disable")
	if err != nil {
		return DBConnection{}, err
	}
	return DBConnection{
		ctx:     ctx,
		queries: sqlc.New(conn),
	}, nil
}

func (conn DBConnection) GetMeals(userId string, date time.Time) ([]sqlc.Meal, error) {
	myUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	meals, err := conn.queries.ListMealsByDate(conn.ctx, sqlc.ListMealsByDateParams{
		UserID:  myUUID,
		ForDate: date,
	})
	if err != nil {
		return nil, err
	}
	return meals, nil
}
