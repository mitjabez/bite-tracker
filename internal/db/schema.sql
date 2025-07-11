-- START: Schema --------------------------------------------------
DROP TABLE IF EXISTS meals;

DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  username text,
  first_name text,
  last_name text,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX idx_username ON users (username);

CREATE TABLE meals (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES users(id) ON DELETE CASCADE,
  meal_type text NOT NULL,
  time_of_meal timestamp NOT NULL,
  description  text NOT NULL,
  hunger_level integer NOT NULL,
  symptoms text[] DEFAULT '{}',
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX idx_meals_time_of_meal ON meals (time_of_meal);
-- END: Schema --------------------------------------------------

-- START: Seed data ---------------------------------------------
-- Insert user with fixed UUID
INSERT INTO users (id, username, first_name, last_name)
VALUES ('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'salsajimmy', 'Salsa', 'Jimmy');

-- Day 1
INSERT INTO meals (user_id, meal_type, time_of_meal, description, hunger_level, symptoms)
VALUES
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'Breakfast', TIMESTAMP '2025-03-01 08:00:00', 'Oatmeal with fruit', 5, ARRAY[]::text[]),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'Lunch',     TIMESTAMP '2025-03-01 13:00:00', 'Grilled chicken sandwich', 6, ARRAY['Bloating']),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'Dinner',    TIMESTAMP '2025-03-01 19:00:00', 'Pasta with tomato sauce', 7, ARRAY['Full']);

-- Day 2
INSERT INTO meals (user_id, meal_type, time_of_meal, description, hunger_level, symptoms)
VALUES
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'Breakfast', TIMESTAMP '2025-03-02 08:00:00', 'Yogurt and granola', 4, ARRAY[]::text[]),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'Lunch',     TIMESTAMP '2025-03-02 12:00:00', 'Turkey salad wrap', 6, ARRAY['Gas']),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'Snack',     TIMESTAMP '2025-03-02 16:00:00', 'Apple slices', 3, ARRAY[]::text[]),
('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'Dinner',    TIMESTAMP '2025-03-02 20:00:00', 'Steamed vegetables and rice', 5, ARRAY['Acid']);

-- END: Seed data ---------------------------------------------
