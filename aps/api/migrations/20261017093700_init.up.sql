CREATE TABLE users (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  telegram_user_id BIGINT NOT NULL UNIQUE,
  username VARCHAR(64),
  name VARCHAR(128),
  created_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE category AS ENUM ('food', 'transport', 'housing', 'entertainment', 'health', 'purchases', 'other');

CREATE TABLE spendings (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id INT NOT NULL,
  amount INT NOT NULL CHECK (amount > 0),
  category category NOT NULL,
  spent_at TIMESTAMPTZ(3) NOT NULL,
  comment TEXT,
  created_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(3),
  CONSTRAINT fk_spendings_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_spendings_user_spent_at ON spendings (user_id, spent_at);
