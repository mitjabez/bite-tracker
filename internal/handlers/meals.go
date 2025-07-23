package handlers

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	db "github.com/mitjabez/bite-tracker/internal/db/init"
	"github.com/mitjabez/bite-tracker/internal/db/sqlc"
	"github.com/mitjabez/bite-tracker/internal/models"
	"github.com/mitjabez/bite-tracker/internal/views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Mealhandler struct {
	dbContext db.DBContext
	// TODO: Move to session
	userId uuid.UUID
}

func NewMealHandler(dbContext db.DBContext, userId string) Mealhandler {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		log.Fatal("Error parsing uuid", userId, err)
	}
	return Mealhandler{dbContext: dbContext, userId: userUUID}
}

func (h Mealhandler) ListMeals(w http.ResponseWriter, r *http.Request) {
	date := dateParam(r)
	currentDate := date.Format("2006-01-02")
	prevDate := date.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	params := sqlc.ListMealsByUsernameAndDateParams{
		UserID:  h.userId,
		ForDate: date,
	}
	meals, err := h.dbContext.Queries.ListMealsByUsernameAndDate(ctx, params)
	if err != nil {
		log.Fatal("Error retrieving meals: ", err)
	}

	mealsView := []models.MealView{}
	for _, m := range meals {
		unusedSymptoms := getUnusedSymptoms(m.Symptoms)
		mealsView = append(mealsView, models.MealView{
			Id:             m.ID.String(),
			MealType:       cases.Title(language.English).String(m.MealTypeID),
			DateOfMeal:     m.TimeOfMeal.Format("2006-01-02"),
			TimeOfMeal:     m.TimeOfMeal.Format("15:04"),
			Description:    m.Description,
			HungerLevel:    m.HungerLevel,
			UsedSymptoms:   m.Symptoms,
			UnusedSymptoms: unusedSymptoms,
		})
	}

	views.Layout(views.Meals(prevDate, nextDate, currentDate, mealsView), "Meal Log").Render(r.Context(), w)
}

func (h Mealhandler) NewMeal(w http.ResponseWriter, r *http.Request) {
	mealView := models.MealView{
		DateOfMeal:  dateParam(r).Format("2006-01-02"),
		TimeOfMeal:  time.Now().Format("15:04"),
		HungerLevel: 4,
	}

	top3Meals, err := h.top3Meals(r.Context())
	if err != nil {
		log.Fatal("Error retrieving top meals for user: ", err)
	}

	views.Layout(views.MealsNew(mealView, map[string]string{}, models.Symptoms, top3Meals), "New Meal").Render(r.Context(), w)
}

func (h Mealhandler) EditMeal(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	userIdParam := r.PathValue("id")
	userUUID, err := uuid.Parse(userIdParam)
	if err != nil {
		log.Fatal("Invalid uuid ", userIdParam)
	}

	meal, err := h.dbContext.Queries.GetMeal(ctx, userUUID)
	if err != nil {
		log.Fatal("Error reading meal: ", err)
	}

	unusedSymptoms := getUnusedSymptoms(meal.Symptoms)
	mealView := models.MealView{
		Id:             userIdParam,
		MealType:       meal.MealTypeID,
		DateOfMeal:     meal.TimeOfMeal.Format("2006-01-02"),
		TimeOfMeal:     meal.TimeOfMeal.Format("15:04"),
		Description:    meal.Description,
		HungerLevel:    meal.HungerLevel,
		UsedSymptoms:   meal.Symptoms,
		UnusedSymptoms: unusedSymptoms,
	}

	top3Meals, err := h.top3Meals(r.Context())
	if err != nil {
		log.Fatal("Error retrieving top meals for user: ", err)
	}
	views.Layout(views.MealsEdit(mealView, map[string]string{}, models.Symptoms, top3Meals), "Edit Meal").Render(r.Context(), w)
}

func (h Mealhandler) HandleMealForm(w http.ResponseWriter, r *http.Request) {
	mealIdParam := r.PathValue("id")
	var mealUUID uuid.UUID
	var err error
	isNewMeal := r.Method == "POST"
	if !isNewMeal {
		mealUUID, err = uuid.Parse(mealIdParam)
		if err != nil {
			log.Fatal("Invalid meal uuid ", mealIdParam)
		}
	}

	dateParam := r.FormValue("date")
	timeParam := r.FormValue("time")
	mealParam := strings.Trim(r.FormValue("meal"), " ")
	hungerParam := r.FormValue("hunger")
	symptoms := r.PostForm["symptoms"]

	errors := map[string]string{}

	if len(mealParam) == 0 {
		errors["meal"] = "Meal is required"
	}

	_, err = time.Parse("2006-01-02", dateParam)
	if err != nil {
		errors["date"] = "Invalid date"
	}

	_, err = time.Parse("15:04", timeParam)
	if err != nil {
		errors["time"] = "Invalid time"
	}
	dateAndTime, err := time.Parse("2006-01-02 15:04", dateParam+" "+timeParam)
	if err != nil {
		errors["time"] = "Invalid date or time"
		errors["date"] = "Invalid date or time"
	}

	if len(mealParam) == 0 {
		errors["meal"] = "Meal is required"
	}

	hungerLevel, err := strconv.Atoi(hungerParam)
	if err != nil || hungerLevel < 1 || hungerLevel > 5 {
		errors["hunger"] = "Invalid hunger level"
	}

	for _, s := range symptoms {
		if !slices.Contains(models.Symptoms, s) {
			errors["symptoms"] = "Symptom " + s + " doesn't exist"
			break
		}
	}

	mealsView := models.MealView{
		DateOfMeal:  dateParam,
		TimeOfMeal:  timeParam,
		Description: mealParam,
		HungerLevel: int32(hungerLevel),
	}

	if len(errors) > 0 {
		top3Meals, err := h.top3Meals(r.Context())
		if err != nil {
			log.Fatal("Error obtaining top 3 meals")
		}
		if isNewMeal {
			views.Layout(views.MealsNew(mealsView, errors, models.Symptoms, top3Meals), "New Meal").Render(r.Context(), w)
		} else {
			views.Layout(views.MealsEdit(mealsView, errors, models.Symptoms, top3Meals), "Edit Meal").Render(r.Context(), w)
		}
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	tx, err := h.dbContext.Pool.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(ctx)

	qtx := h.dbContext.Queries.WithTx(tx)

	mealType := resolveMealType(dateAndTime)
	if isNewMeal {
		_, err = qtx.CreateMeal(ctx, sqlc.CreateMealParams{
			MealTypeID:  mealType,
			TimeOfMeal:  dateAndTime,
			Description: mealParam,
			HungerLevel: int32(hungerLevel),
			Symptoms:    symptoms,
			UserID:      h.userId,
		})
	} else {
		err = qtx.UpdateMeal(ctx, sqlc.UpdateMealParams{
			ID:          mealUUID,
			MealTypeID:  mealType,
			TimeOfMeal:  dateAndTime,
			Description: mealParam,
			HungerLevel: int32(hungerLevel),
			Symptoms:    symptoms,
			UpdatedAt:   time.Now(),
		})
	}
	if err != nil {
		log.Fatal("Cannot create or update meal: ", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("Cannot commit meals transaction: ", err)
	}

	h.ListMeals(w, r)
}

func (h Mealhandler) top3Meals(parentContext context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(parentContext, 2*time.Second)
	defer cancel()
	top3MealsResult, err := h.dbContext.Queries.Top3Meals(ctx, sqlc.Top3MealsParams{
		UserID:     h.userId,
		MealTypeID: resolveMealType(time.Now()),
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

func resolveMealType(time time.Time) string {
	hour := time.Hour()
	switch {
	case hour < 9:
		return "breakfast"
	case hour < 11:
		return "brunch"
	case hour < 15:
		return "lunch"
	default:
		return "dinner"
	}
}

func dateParam(r *http.Request) time.Time {
	dateQuery := r.FormValue("date")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if dateQuery != "" {
		parsedDate, err := time.Parse("2006-01-02", dateQuery)
		if err != nil {
			log.Println("WARNING: Error parsing date", dateQuery)
		} else {
			date = parsedDate
		}
	}
	return date
}

// TODO: Optimize
func getUnusedSymptoms(activeSymptoms []string) []string {
	usedSymptomNames := map[string]bool{}
	for _, s := range activeSymptoms {
		usedSymptomNames[s] = true
	}

	unusedSymptoms := []string{}
	for _, s := range models.Symptoms {
		if !usedSymptomNames[s] {
			unusedSymptoms = append(unusedSymptoms, s)
		}
	}
	return unusedSymptoms
}
