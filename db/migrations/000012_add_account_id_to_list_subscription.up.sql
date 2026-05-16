ALTER TABLE list_subscriptions ADD COLUMN account_id BIGINT NOT NULL;
ALTER TABLE list_subscriptions ADD CONSTRAINT fk_list_subscriptions_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX idx_list_subscriptions_account_id_list_id_contact_id ON list_subscriptions (account_id, list_id, contact_id);