CREATE TABLE configurations (
  id                      BIGSERIAL PRIMARY KEY,
  account_id              BIGINT NOT NULL,
  is_default               BOOLEAN DEFAULT FALSE, -- this defines if this is default configuration for account
  config_type              INT NOT NULL, -- 0 = smtp, 1 = web_app
  settings                JSONB DEFAULT '{}'::JSONB, -- this defines configuration
  created_at              TIMESTAMP DEFAULT now(),
  updated_at              TIMESTAMP DEFAULT now(),

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_configurations_account_id ON configurations (account_id);

CREATE UNIQUE INDEX idx_configurations_account_id_type_default ON configurations (account_id, config_type, is_default);