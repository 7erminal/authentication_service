-- MySQL migration for currently used persisted models in authentication_service
-- Excludes non-persistent DTOs (for example AuthenticationDTO).

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS customer_credentials (
  customer_credential_id BIGINT NOT NULL AUTO_INCREMENT,
  customer_id BIGINT NOT NULL,
  username VARCHAR(255) NOT NULL DEFAULT '',
  password VARCHAR(255) NOT NULL DEFAULT '',
  pin VARCHAR(10) NOT NULL DEFAULT '',
  date_created DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  created_by INT NOT NULL DEFAULT 0,
  modified_by INT NOT NULL DEFAULT 0,
  active INT NOT NULL DEFAULT 0,
  PRIMARY KEY (customer_credential_id),
  KEY idx_customer_credentials_customer_id (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS access_tokens (
  access_token_id BIGINT NOT NULL AUTO_INCREMENT,
  token VARCHAR(255) NOT NULL DEFAULT '',
  user_id BIGINT NOT NULL,
  date_created DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT 0,
  ip_address VARCHAR(80) NOT NULL DEFAULT '',
  last_used_at DATETIME NULL,
  PRIMARY KEY (access_token_id),
  KEY idx_access_tokens_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS refresh_tokens (
  refresh_token_id BIGINT NOT NULL AUTO_INCREMENT,
  token VARCHAR(255) NOT NULL,
  user_id BIGINT NOT NULL,
  access_token_id BIGINT NULL,
  date_created DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT 0,
  i_p_address VARCHAR(45) NULL,
  user_agent VARCHAR(255) NULL,
  last_used_at DATETIME NULL,
  PRIMARY KEY (refresh_token_id),
  UNIQUE KEY uq_refresh_tokens_token (token),
  KEY idx_refresh_tokens_user_id (user_id),
  KEY idx_refresh_tokens_access_token_id (access_token_id),
  CONSTRAINT fk_refresh_tokens_access_token FOREIGN KEY (access_token_id) REFERENCES access_tokens(access_token_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS customer_access_tokens (
  customer_access_token_id BIGINT NOT NULL AUTO_INCREMENT,
  token VARCHAR(255) NOT NULL DEFAULT '',
  customer_id BIGINT NOT NULL,
  date_created DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT 0,
  ip_address VARCHAR(80) NOT NULL DEFAULT '',
  last_used_at DATETIME NOT NULL,
  PRIMARY KEY (customer_access_token_id),
  KEY idx_customer_access_tokens_customer_id (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS customer_refresh_tokens (
  refresh_token_id BIGINT NOT NULL AUTO_INCREMENT,
  token VARCHAR(255) NOT NULL,
  customer_id BIGINT NOT NULL,
  access_token_id BIGINT NULL,
  date_created DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT 0,
  i_p_address VARCHAR(45) NULL,
  user_agent VARCHAR(255) NULL,
  last_used_at DATETIME NULL,
  PRIMARY KEY (refresh_token_id),
  UNIQUE KEY uq_customer_refresh_tokens_token (token),
  KEY idx_customer_refresh_tokens_customer_id (customer_id),
  KEY idx_customer_refresh_tokens_access_token_id (access_token_id),
  CONSTRAINT fk_customer_refresh_tokens_access_token FOREIGN KEY (access_token_id) REFERENCES customer_access_tokens(customer_access_token_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_tokens (
  user_token_id BIGINT NOT NULL AUTO_INCREMENT,
  token VARCHAR(255) NOT NULL DEFAULT '',
  nonce VARCHAR(255) NOT NULL DEFAULT '',
  expiry_date DATETIME NOT NULL,
  date_created DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  created_by INT NOT NULL DEFAULT 0,
  modified_by INT NOT NULL DEFAULT 0,
  active INT NOT NULL DEFAULT 0,
  PRIMARY KEY (user_token_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_otps (
  user_otp_id BIGINT NOT NULL AUTO_INCREMENT,
  code VARCHAR(128) NOT NULL DEFAULT '',
  user_id BIGINT NOT NULL,
  status INT NOT NULL DEFAULT 0,
  date_created DATETIME NOT NULL,
  date_generated DATETIME NOT NULL,
  expiry_date DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  created_by INT NOT NULL DEFAULT 0,
  modified_by INT NOT NULL DEFAULT 0,
  active INT NOT NULL DEFAULT 0,
  PRIMARY KEY (user_otp_id),
  KEY idx_user_otps_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS activation_codes (
  activation_code_id BIGINT NOT NULL AUTO_INCREMENT,
  code VARCHAR(80) NOT NULL DEFAULT '',
  number VARCHAR(80) NOT NULL DEFAULT '',
  expiry_date DATETIME NOT NULL,
  date_created DATETIME NOT NULL,
  date_modified DATETIME NOT NULL,
  created_by INT NOT NULL DEFAULT 0,
  modified_by INT NOT NULL DEFAULT 0,
  active INT NOT NULL DEFAULT 0,
  PRIMARY KEY (activation_code_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
