package mealservice

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func New(connectionString string) (DBConnection, error) {
	// TODO: Handle timeout
	ctx := context.Background() //context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return DBConnection{}, err
	}
	return DBConnection{
		ctx:     ctx,
		queries: sqlc.New(conn),
	}, nil
}

func (conn DBConnection) GetMeals(username string, date time.Time) ([]sqlc.Meal, error) {
	meals, err := conn.queries.ListMealsByUsernameAndDate(conn.ctx, sqlc.ListMealsByUsernameAndDateParams{
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
		ForDate: date,
	})
	if err != nil {
		return nil, err
	}
	return meals, nil
}
