CREATE TABLE email_notifications (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  subject VARCHAR(500) NOT NULL CHECK (subject <> ''),
  title VARCHAR(300) NOT NULL CHECK (title <> ''),
  notification_type INTEGER NOT NULL CHECK (notification_type IN (0, 1)), -- 0 = transactional, 1 = campaign(sent to bulk contacts)
  content_id BIGINT,

  status INTEGER NOT NULL DEFAULT 0 CHECK (status IN (0, 1, 2, 3)), -- [0 - pending, 1 - enqueued, 2 - sent, 3 - failed]
  sent_at TIMESTAMP, -- delivered time
  created_at TIMESTAMP DEFAULT now(),
  send_at TIMESTAMP, -- time for triggering delivery from system

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_email_notifications_account_id ON email_notifications (account_id);
CREATE INDEX idx_email_notifications_account_id_status ON email_notifications (account_id, status);
