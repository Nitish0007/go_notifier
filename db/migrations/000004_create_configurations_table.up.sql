CREATE TABLE configurations (
  id                      SERIAL PRIMARY KEY,
  account_id              INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  default_configuration   BOOLEAN DEFAULT FALSE, -- this defines if this is default configuration for account
  config_type             TEXT NOT NULL, -- this defines configuration type like : logger(if user want logs on his logger), provider like email/sms etc
  configuration_data      JSONB DEFAULT '{}'::JSONB, -- this defines configuration
  created_at              TIMESTAMP DEFAULT now(),
  updated_at              TIMESTAMP DEFAULT now(),

  CONSTRAINT config_type_format CHECK (config_type ~ '^[a-zA-Z0-9_]+$')
  -- This regex allows alphanumeric characters and underscores, ensuring a valid config type format
);

CREATE INDEX idx_configurations_account_id ON configurations (account_id);