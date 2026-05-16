ALTER TABLE email_notifications
  DROP COLUMN IF EXISTS reply_to_email,
  DROP COLUMN IF EXISTS from_email,
  DROP COLUMN IF EXISTS from_name;
