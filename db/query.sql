-- name: ListMeals :many
SELECT * FROM meals
WHERE time_of_meal = $1
ORDER BY time_of_meal;

-- name: CreateMeal :one
INSERT INTO meals (
  user_id,
  meal_type,
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
  SET meal_type = $2,
  time_of_meal = $3,
  description = $4,
  hunger_level = $5,
  symptoms = $6,
  updated_at = $7
WHERE id = $1;

