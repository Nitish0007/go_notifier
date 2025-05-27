CREATE TABLE api_keys (
  id                  SERIAL PRIMARY KEY,
  account_id          INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  key                 TEXT NOT NULL UNIQUE,
  created_at          TIMESTAMP DEFAULT now(),
  updated_at          TIMESTAMP DEFAULT now(),
  
  CONSTRAINT key_format CHECK (key ~ '^[a-zA-Z0-9]{32}$')
);