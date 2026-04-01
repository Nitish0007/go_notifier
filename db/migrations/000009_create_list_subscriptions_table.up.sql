CREATE TABLE list_subscriptions (
  id BIGSERIAL PRIMARY KEY,
  list_id BIGINT NOT NULL,
  contact_id BIGINT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),

  FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE,
  FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);

CREATE INDEX idx_list_subscriptions_list_id ON list_subscriptions (list_id);
CREATE INDEX idx_list_subscriptions_contact_id ON list_subscriptions (contact_id);
CREATE UNIQUE INDEX idx_list_subscriptions_list_id_contact_id ON list_subscriptions (list_id, contact_id);