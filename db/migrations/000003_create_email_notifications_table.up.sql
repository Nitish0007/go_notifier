CREATE TABLE email_notifications (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  subject VARCHAR(500) NOT NULL CHECK (subject <> ''),
  title VARCHAR(300) NOT NULL CHECK (title <> ''),
  notification_type INTEGER NOT NULL CHECK (notification_type IN (0, 1)),
  -- 0 = transactional, 1 = campaign

  content_id BIGINT NOT NULL,

  -- Matches Go emailnotification.EmailNotificationStatus (iota 0..5)
  -- 0 Trans, 1 Draft, 2 Scheduled, 3 Enqueued, 4 Sent, 5 Failed
  status INTEGER NOT NULL DEFAULT 0 CHECK (status IN (0, 1, 2, 3, 4, 5)),

  send_at TIMESTAMP,
  sent_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_email_notifications_account_id ON email_notifications (account_id);
CREATE INDEX idx_email_notifications_account_id_status ON email_notifications (account_id, status);
