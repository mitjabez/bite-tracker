# Go + ah templ Meals App Conventions

## Project structure

```
/internal
├── handlers
│   └── meal.go        // package handlers
├── models
│   └── meal.go        // package models
├── templates
│   ├── layout.templ
│   ├── meals.templ
│   ├── meals_new.templ
│   ├── meal.templ
│   └── meal_edit.templ
/cmd
└── main.go
```

## Package names

- `handlers/meal.go` → package handlers
- `models/meal.go` → package models

## Structs

- `MealHandler` (in handlers): holds DB & context, handles `/meals` HTTP routes
- `Meal` (in models): data struct with fields like ID, Name, Date

---

## Handler functions (in handlers/meal.go)

```go
func (h MealHandler) ListMeals(w http.ResponseWriter, r *http.Request)
func (h MealHandler) ShowMeal(w http.ResponseWriter, r *http.Request)
func (h MealHandler) NewMeal(w http.ResponseWriter, r *http.Request)
func (h MealHandler) CreateMeal(w http.ResponseWriter, r *http.Request)
func (h MealHandler) EditMeal(w http.ResponseWriter, r *http.Request)
func (h MealHandler) UpdateMeal(w http.ResponseWriter, r *http.Request)
func (h MealHandler) DeleteMeal(w http.ResponseWriter, r *http.Request)
```

## Model functions (in models/meal.go)

```go
func GetAllMeals() ([]Meal, error)
func GetMealByID(id int) (Meal, error)
func CreateMeal(m *Meal) error
func UpdateMeal(m *Meal) error
func DeleteMeal(id int) error
```

## Templates

- `layout.templ` → shared HTML layout
- `meals.templ` → list of meals page
- `meals_new.templ` → form to add a meal
- `meal.templ` → show single meal
- `meal_edit.templ` → form to edit meal

---

## Summary of conventions

- Group by resource, not by individual page.
- Keep all `/meals` HTTP handlers in `meal.go`.
- Use plural names for collections (`meals.templ`), singular for single items (`meal.templ`).
- Do all logic in Go, pass prepared data to templates, keep templates purely for rendering.
