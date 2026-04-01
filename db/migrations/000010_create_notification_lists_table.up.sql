CREATE TABLE email_notification_lists (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  list_id BIGINT NOT NULL,
  notification_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (list_id) REFERENCES lists(id),
  FOREIGN KEY (notification_id) REFERENCES email_notifications(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_email_notification_lists_account_id_list_id_notification_id ON email_notification_lists (account_id, list_id, notification_id);