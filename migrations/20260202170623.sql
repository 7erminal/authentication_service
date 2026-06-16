-- Create "customer_refresh_tokens" table
CREATE TABLE customer_refresh_tokens (
  refresh_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL DEFAULT '',
  customer_id BIGINT NOT NULL,
  access_token_id BIGINT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  i_p_address VARCHAR(45) NULL,
  user_agent VARCHAR(255) NULL,
  last_used_at TIMESTAMP NULL,
  UNIQUE (token)
);
