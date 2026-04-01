CREATE TABLE lists (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,
  contacts_count BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_lists_account_id_name ON lists (account_id, name);