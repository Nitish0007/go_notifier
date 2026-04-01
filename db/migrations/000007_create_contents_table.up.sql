CREATE TABLE contents (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  body TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  
  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);