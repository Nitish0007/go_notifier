CREATE TABLE notification_batches (
  id UUID DEFAULT gen_random_uuid(),
  account_id INT NOT NULL,
  count int DEFAULT 0,
  successful_count int DEFAULT 0,
  failed_count int DEFAULT 0,
  channel INTEGER NOT NULL CHECK (channel IN (0, 1, 2)), -- e.g., 0=email, 1=sms, 2=in_app
  status INTEGER NOT NULL DEFAULT 0 CHECK (status IN (0, 1, 2, 3)), -- [0 - pending, 1 - enqueued, 2 - sent, 3 - failed]
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  completed_at TIMESTAMP DEFAULT NULL,

  PRIMARY KEY (id),
  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_notification_batches_account_id 
  ON notification_batches (account_id);

CREATE INDEX idx_notification_batches_account_id_channel_status 
  ON notification_batches (account_id, channel, status);

CREATE INDEX idx_notification_batches_account_id_created_at_completed_at 
  ON notification_batches (account_id, created_at, completed_at);