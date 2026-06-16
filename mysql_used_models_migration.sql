-- Authentication Service migration scoped to DB-backed models
-- currently present in this project.
--
-- Included model tables:
--   access_tokens
--   refresh_tokens
--   customer_access_tokens
--   customer_refresh_tokens
--   customer_credentials
--   user_tokens
--   user_otps
--   activation_codes
--
-- Note: AuthenticationDTO is intentionally excluded because it is a request DTO,
-- not a persistent table model.

CREATE TABLE IF NOT EXISTS customer_credentials (
  customer_credential_id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL,
  username VARCHAR(255) NOT NULL DEFAULT '',
  password VARCHAR(255) NOT NULL DEFAULT '',
  pin VARCHAR(10) NOT NULL DEFAULT '',
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_customer_credentials_customer_id ON customer_credentials(customer_id);

CREATE TABLE IF NOT EXISTS access_tokens (
  access_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL DEFAULT '',
  user_id BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  ip_address VARCHAR(80) NOT NULL DEFAULT '',
  last_used_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_access_tokens_user_id ON access_tokens(user_id);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  refresh_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL,
  user_id BIGINT NOT NULL,
  access_token_id BIGINT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  i_p_address VARCHAR(45) NULL,
  user_agent VARCHAR(255) NULL,
  last_used_at TIMESTAMP NULL,
  CONSTRAINT uq_refresh_tokens_token UNIQUE (token),
  CONSTRAINT fk_refresh_tokens_access_token FOREIGN KEY (access_token_id) REFERENCES access_tokens(access_token_id)
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_access_token_id ON refresh_tokens(access_token_id);

CREATE TABLE IF NOT EXISTS customer_access_tokens (
  customer_access_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL DEFAULT '',
  customer_id BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  ip_address VARCHAR(80) NOT NULL DEFAULT '',
  last_used_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_customer_access_tokens_customer_id ON customer_access_tokens(customer_id);

CREATE TABLE IF NOT EXISTS customer_refresh_tokens (
  refresh_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL,
  customer_id BIGINT NOT NULL,
  access_token_id BIGINT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  i_p_address VARCHAR(45) NULL,
  user_agent VARCHAR(255) NULL,
  last_used_at TIMESTAMP NULL,
  CONSTRAINT uq_customer_refresh_tokens_token UNIQUE (token),
  CONSTRAINT fk_customer_refresh_tokens_access_token FOREIGN KEY (access_token_id) REFERENCES customer_access_tokens(customer_access_token_id)
);
CREATE INDEX IF NOT EXISTS idx_customer_refresh_tokens_customer_id ON customer_refresh_tokens(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_refresh_tokens_access_token_id ON customer_refresh_tokens(access_token_id);

CREATE TABLE IF NOT EXISTS user_tokens (
  user_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL DEFAULT '',
  nonce VARCHAR(255) NOT NULL DEFAULT '',
  expiry_date TIMESTAMP NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_otps (
  user_otp_id BIGSERIAL PRIMARY KEY,
  code VARCHAR(128) NOT NULL DEFAULT '',
  user_id BIGINT NOT NULL,
  status INTEGER NOT NULL DEFAULT 0,
  date_created TIMESTAMP NOT NULL,
  date_generated TIMESTAMP NOT NULL,
  expiry_date TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_user_otps_user_id ON user_otps(user_id);

CREATE TABLE IF NOT EXISTS activation_codes (
  activation_code_id BIGSERIAL PRIMARY KEY,
  code VARCHAR(80) NOT NULL DEFAULT '',
  number VARCHAR(80) NOT NULL DEFAULT '',
  expiry_date TIMESTAMP NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);
