ALTER TABLE contents ADD COLUMN name VARCHAR(100) NOT NULL;

CREATE UNIQUE INDEX idx_content_account_id_name ON contents (account_id, name);