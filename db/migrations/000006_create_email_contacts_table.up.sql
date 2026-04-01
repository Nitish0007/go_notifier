CREATE TABLE email_contacts (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL,
  contact_id BIGINT NOT NULL,
  email VARCHAR(255) NOT NULL CHECK (email ~ '^[^@]+@[^@]+\.[^@]+$'),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),

  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_email_contacts_account_id ON email_contacts (account_id);
CREATE INDEX idx_email_contacts_account_id_contact_id ON email_contacts (account_id, contact_id);
CREATE UNIQUE INDEX idx_email_contacts_email_account_id ON email_contacts (email, account_id);