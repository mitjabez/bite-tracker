CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text NOT NULL,
  full_name text NOT NULL,
	password_hash text NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

INSERT INTO users (id, email, full_name, password_hash)
VALUES ('f41ad27a-881d-4f7f-a908-f16a26ce7b78', 'sj@dot.com', 'Salsa Jimmy', '$2a$12$F22j/9fE8wI2nfjFADc/reQgm/TpKAxUWIyPhzZybV3GuvZP49rtu');
