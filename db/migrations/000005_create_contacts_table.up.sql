CREATE TABLE contacts (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_contacts_account_id ON contacts (account_id);
CREATE UNIQUE INDEX idx_contacts_account_id_uuid ON contacts (account_id, uuid);