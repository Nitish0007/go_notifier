CREATE TABLE notification_batch_errors (
  id UUID DEFAULT gen_random_uuid(),
  batch_id UUID NOT NULL,
  error TEXT NOT NULL, -- error message
  payload JSONB DEFAULT '{}'::JSONB NOT NULL, -- payload of the notification that caused error
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (batch_id) REFERENCES notification_batches(id) ON DELETE CASCADE
);