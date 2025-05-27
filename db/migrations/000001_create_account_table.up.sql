CREATE TABLE accounts (
  id                  SERIAL PRIMARY KEY,
  email               TEXT NOT NULL,
  encrypted_password  TEXT NOT NULL,
  first_name          TEXT NOT NULL,
  last_name           TEXT NOT NULL,
  is_active           BOOLEAN DEFAULT TRUE,
  created_at          TIMESTAMP DEFAULT now(),
  updated_at          TIMESTAMP DEFAULT now(),
  
  CONSTRAINT email_unique UNIQUE (email),
  CONSTRAINT email_format CHECK (email ~ '^[^@]+@[^@]+\.[^@]+$')
)