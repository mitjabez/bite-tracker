-- name: ListMealsByUsernameAndDate :many
SELECT * FROM meals
WHERE user_id = @user_id AND
	time_of_meal > @for_date AND time_of_meal < ( (@for_date) + interval '1 day' )
ORDER BY time_of_meal;

-- name: CreateMeal :one
INSERT INTO meals (
  user_id,
  meal_type_id,
  time_of_meal,
  description,
  hunger_level,
  symptoms
) VALUES (
  $1,$2,$3,$4,$5,$6
)
RETURNING *;

-- name: UpdateMeal :exec
UPDATE meals
  SET meal_type_id = $2,
  time_of_meal = $3,
  description = $4,
  hunger_level = $5,
  symptoms = $6,
  updated_at = $7
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM user
WHERE u.username = @username::text
LIMIT 1;

-- name: Top3Meals :many
SELECT description, times_used FROM meals_catalog
WHERE user_id = @user_id AND meal_type_id = @meal_type_id
ORDER BY times_used DESC
LIMIT 3;

-- name: UpdateMealsCatalog :exec
INSERT INTO meals_catalog (
	user_id,
	description,
  meal_type_id
) VALUES (
	@user_id,
	@description,
	@meal_type_id
) ON CONFLICT(user_id, description, meal_type_id) DO UPDATE
SET times_used = meals_catalog.times_used + 1;
