package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
	"github.com/mitjabez/bite-tracker/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MealRepo struct {
	DBContext db.DBContext
}

func (r *MealRepo) ListMeals(ctx context.Context, userId uuid.UUID, date time.Time) ([]model.MealView, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	params := sqlc.ListMealsByUsernameAndDateParams{
		UserID:  userId,
		ForDate: date,
	}
	meals, err := r.DBContext.Queries.ListMealsByUsernameAndDate(ctx, params)
	if err != nil {
		return []model.MealView{}, err
	}

	mealsView := []model.MealView{}
	for _, m := range meals {
		mealsView = append(mealsView, mealToMealView(m))
	}
	return mealsView, nil
}

func (r *MealRepo) GetMeal(ctx context.Context, userId uuid.UUID) (model.MealView, error) {
	meal, err := r.DBContext.Queries.GetMeal(ctx, userId)
	if err != nil {
		return model.MealView{}, nil
	}

	return mealToMealView(meal), nil
}

func (r *MealRepo) CreateMeal(ctx context.Context, userId uuid.UUID, mealView model.MealView) error {
	return r.createOrUpdateMeal(ctx, true, userId, uuid.Nil, mealView)
}

func (r *MealRepo) UpdateMeal(ctx context.Context, userId uuid.UUID, mealId uuid.UUID, mealView model.MealView) error {
	return r.createOrUpdateMeal(ctx, false, userId, mealId, mealView)
}

func (r *MealRepo) createOrUpdateMeal(ctx context.Context, isNewMeal bool, userId uuid.UUID, mealId uuid.UUID, mealView model.MealView) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	tx, err := r.DBContext.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.DBContext.Queries.WithTx(tx)

	if isNewMeal {
		_, err = qtx.CreateMeal(ctx, sqlc.CreateMealParams{
			UserID:      userId,
			MealTypeID:  mealView.MealType,
			TimeOfMeal:  mealView.TimeOfMeal,
			Description: mealView.Description,
			HungerLevel: mealView.HungerLevel,
			Symptoms:    mealView.Symptoms,
		})
	} else {
		err = qtx.UpdateMeal(ctx, sqlc.UpdateMealParams{
			ID:          mealId,
			MealTypeID:  mealView.MealType,
			TimeOfMeal:  mealView.TimeOfMeal,
			Description: mealView.Description,
			HungerLevel: mealView.HungerLevel,
			Symptoms:    mealView.Symptoms,
			UpdatedAt:   time.Now(),
		})
	}
	if err != nil {
		return err
	}
	err = qtx.UpdateMealsCatalog(ctx, sqlc.UpdateMealsCatalogParams{
		UserID:      userId,
		Description: mealView.Description,
		MealTypeID:  mealView.MealType,
	})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *MealRepo) Top3Meals(ctx context.Context, userId uuid.UUID) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	top3MealsResult, err := r.DBContext.Queries.Top3Meals(ctx, sqlc.Top3MealsParams{
		UserID:     userId,
		MealTypeID: model.ResolveMealType(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	top3Meals := []string{}
	for _, m := range top3MealsResult {
		top3Meals = append(top3Meals, m.Description)
	}
	return top3Meals, nil
}

func mealToMealView(m sqlc.Meal) model.MealView {
	return model.MealView{
		Id:          m.ID.String(),
		MealType:    cases.Title(language.English).String(m.MealTypeID),
		TimeOfMeal:  m.TimeOfMeal,
		Description: m.Description,
		HungerLevel: m.HungerLevel,
		Symptoms:    m.Symptoms,
	}
}
