CREATE TABLE email_sends (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  notification_id BIGINT NOT NULL,
  contact_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (notification_id) REFERENCES email_notifications(id) ON DELETE CASCADE,
  FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);

CREATE INDEX idx_email_sends_account_id_notification_id_contact_id ON email_sends (account_id, notification_id, contact_id);
CREATE INDEX idx_email_sends_account_id_contact_id ON email_sends (account_id, contact_id);
