CREATE TABLE notifications (
  id UUID DEFAULT gen_random_uuid(),
  account_id INT NOT NULL,
  channel INTEGER NOT NULL CHECK (channel IN (0, 1, 2)), -- e.g., 0=email, 1=sms, 2=in_app
  recipient VARCHAR(255) NOT NULL,
  subject VARCHAR(500),
  body TEXT,
  html_body TEXT,
  job_id UUID,
  metadata JSONB DEFAULT '{}'::JSONB,

  status INTEGER NOT NULL DEFAULT 0 CHECK (status IN (0, 1, 2, 3)), -- [0 - pending, 1 - enqueued, 2 - sent, 3 - failed]
  sent_at TIMESTAMP, -- delivered time
  created_at TIMESTAMP DEFAULT NOW(),
  send_at TIMESTAMP, -- time for triggering delivery from system
  error_message TEXT,

  PRIMARY KEY (id, channel),
  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) PARTITION BY LIST(channel);

-- Create partitions
CREATE TABLE notifications_email PARTITION OF notifications FOR VALUES IN (0);
CREATE TABLE notifications_sms PARTITION OF notifications FOR VALUES IN (1);
CREATE TABLE notifications_in_app PARTITION OF notifications FOR VALUES IN (2);

-- Indexes on parent table
CREATE INDEX idx_notifications_account_id ON notifications (account_id);

-- Index on channel specific table (on partitions)
CREATE INDEX idx_notifications_email_status ON notifications_email (status);
CREATE INDEX idx_notifications_sms_status ON notifications_sms (status);
CREATE INDEX idx_notifications_in_app_status ON notifications_in_app (status);

-- Partition-specific constraints
ALTER TABLE notifications_email 
ADD CONSTRAINT email_recipient_format 
CHECK (recipient ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');

ALTER TABLE notifications_sms 
ADD CONSTRAINT sms_recipient_format 
CHECK (recipient ~ '^\+?[1-9]\d{1,14}$');