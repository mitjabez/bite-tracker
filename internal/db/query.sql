-- name: ListMealsByUsernameAndDate :many
SELECT m.* FROM meals m
JOIN users u ON m.user_id = u.id
WHERE u.username = @username::text AND
	time_of_meal > @for_date::timestamp AND time_of_meal < ( (@for_date::timestamp) + interval '1 day' )
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

