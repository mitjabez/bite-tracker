-- START: Schema --------------------------------------------------
DROP TABLE IF EXISTS meals;
DROP TABLE IF EXISTS meals_catalog;
DROP TABLE IF EXISTS meal_types;
DROP TABLE IF EXISTS users CASCADE;


CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  username text NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uniq_users_username ON users (username);

CREATE TABLE meal_types (
  id text PRIMARY KEY,
	start_time time NOT NULL,
	end_time time NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE meals (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  meal_type_id text NOT NULL REFERENCES meal_types(id),
  time_of_meal timestamp NOT NULL,
  description text NOT NULL,
  hunger_level integer NOT NULL CHECK (hunger_level >=1 AND hunger_level <=5),
  symptoms text[] DEFAULT '{}',
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX idx_meals_time_of_meal ON meals (time_of_meal);

CREATE TABLE meals_catalog (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL REFERENCES users(id),
	times_used integer NOT NULL DEFAULT 1,
	description text NOT NULL,
  meal_type_id text NOT NULL REFERENCES meal_types(id),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_meals_catalog_user_mealtype_usage ON meals_catalog (user_id,meal_type_id, times_used);

-- END: Schema --------------------------------------------------

-- START: Seed data ---------------------------------------------
-- Insert user with fixed UUID
INSERT INTO users (id, username, first_name, last_name)
VALUES ('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'salsajimmy', 'Salsa', 'Jimmy');

INSERT INTO meal_types (id, start_time, end_time)
VALUES
('breakfast', '00:00:00', '09:00:00'),
('brunch', '09:00:00', '11:00:00'),
('lunch', '11:00:00', '15:00:00'),
('dinner', '15:00:00', '23:59:59');

-- Day 1
INSERT INTO meals (user_id, meal_type_id, time_of_meal, description, hunger_level, symptoms)
VALUES
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'breakfast', TIMESTAMP '2025-03-01 08:00:00', 'Oatmeal with fruit', 5, ARRAY[]::text[]),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'lunch',     TIMESTAMP '2025-03-01 13:00:00', 'Grilled chicken sandwich', 4, ARRAY['bloating']),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'dinner',    TIMESTAMP '2025-03-01 19:00:00', 'Pasta with tomato sauce', 3, ARRAY['full']);

-- Day 2
INSERT INTO meals (user_id, meal_type_id, time_of_meal, description, hunger_level, symptoms)
VALUES
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'breakfast', TIMESTAMP '2025-03-02 08:00:00', 'Yogurt and granola', 4, ARRAY[]::text[]),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'lunch',     TIMESTAMP '2025-03-02 12:00:00', 'Turkey salad wrap', 2, ARRAY['gas']),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'dinner',     TIMESTAMP '2025-03-02 16:00:00', 'Apple slices', 1, ARRAY[]::text[]),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'dinner',    TIMESTAMP '2025-03-02 20:00:00', 'Steamed vegetables and rice', 4, ARRAY['acid']::text[]);

-- END: Seed data ---------------------------------------------
