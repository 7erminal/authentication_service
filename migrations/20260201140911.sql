-- Create "access_tokens" table
CREATE TABLE access_tokens (
  access_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL DEFAULT '',
  user_id BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  i_p_address VARCHAR(80) NOT NULL DEFAULT '',
  last_used_at TIMESTAMP NOT NULL
);

-- Create "activation_codes" table
CREATE TABLE activation_codes (
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

-- Create "auth_users" table
CREATE TABLE auth_users (
  user_id BIGSERIAL PRIMARY KEY,
  user_type INTEGER NOT NULL DEFAULT 0,
  user_details_id BIGINT NULL,
  image_path VARCHAR(200) NULL,
  full_name VARCHAR(255) NOT NULL DEFAULT '',
  username VARCHAR(255) NOT NULL DEFAULT '',
  password VARCHAR(255) NOT NULL DEFAULT '',
  email VARCHAR(255) NOT NULL DEFAULT '',
  phone_number VARCHAR(255) NOT NULL DEFAULT '',
  gender VARCHAR(10) NOT NULL DEFAULT '',
  dob TIMESTAMP NOT NULL,
  address VARCHAR(255) NOT NULL DEFAULT '',
  id_type VARCHAR(5) NOT NULL DEFAULT '',
  id_number VARCHAR(100) NOT NULL DEFAULT '',
  marital_status VARCHAR(255) NOT NULL DEFAULT '',
  role BIGINT NULL,
  active INTEGER NOT NULL DEFAULT 0,
  is_verified BOOLEAN NOT NULL DEFAULT FALSE,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0
);

-- Create "authentication_d_t_o" table
CREATE TABLE authentication_d_t_o (
  username VARCHAR(255) NOT NULL DEFAULT '',
  password VARCHAR(255) NOT NULL DEFAULT ''
);

-- Create "branches" table
CREATE TABLE branches (
  branch_id BIGSERIAL PRIMARY KEY,
  branch VARCHAR(80) NOT NULL DEFAULT '',
  country_id BIGINT NOT NULL,
  location VARCHAR(255) NOT NULL DEFAULT '',
  phone_number VARCHAR(255) NOT NULL DEFAULT '',
  active INTEGER NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NULL,
  modified_by INTEGER NULL
);

-- Create "countries" table
CREATE TABLE countries (
  country_id BIGSERIAL PRIMARY KEY,
  country VARCHAR(255) NOT NULL DEFAULT '',
  description VARCHAR(500) NOT NULL DEFAULT '',
  country_code VARCHAR(20) NOT NULL DEFAULT '',
  default_currency BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0
);

-- Create "currencies" table
CREATE TABLE currencies (
  currency_id BIGSERIAL PRIMARY KEY,
  symbol VARCHAR(20) NOT NULL DEFAULT '',
  currency VARCHAR(50) NOT NULL DEFAULT '',
  active INTEGER NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NULL,
  modified_by INTEGER NULL
);

-- Create "customer_access_tokens" table
CREATE TABLE customer_access_tokens (
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

-- Create "customer_categories" table
CREATE TABLE customer_categories (
  customer_category_id BIGSERIAL PRIMARY KEY,
  category VARCHAR(100) NOT NULL DEFAULT '',
  description VARCHAR(255) NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);

-- Create "customer_credentials" table
CREATE TABLE customer_credentials (
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

-- Create "customer_emergency_contacts" table
CREATE TABLE customer_emergency_contacts (
  customer_emergency_contact_id BIGSERIAL PRIMARY KEY,
  name VARCHAR(120) NOT NULL DEFAULT '',
  contact VARCHAR(50) NOT NULL DEFAULT '',
  customer_id BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0
);

-- Create "customer_guarantors" table
CREATE TABLE customer_guarantors (
  customer_guarantor_id BIGSERIAL PRIMARY KEY,
  name VARCHAR(120) NOT NULL DEFAULT '',
  contact VARCHAR(50) NOT NULL DEFAULT '',
  customer_id BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0
);

-- Create "customers" table
CREATE TABLE customers (
  customer_id BIGSERIAL PRIMARY KEY,
  customer_number VARCHAR(255) NOT NULL DEFAULT '',
  full_name VARCHAR(255) NOT NULL DEFAULT '',
  image_path VARCHAR(255) NOT NULL DEFAULT '',
  email VARCHAR(255) NULL,
  phone_number VARCHAR(255) NULL,
  location VARCHAR(255) NULL,
  identification_type_id BIGINT NULL,
  identification_number VARCHAR(255) NULL,
  branch BIGINT NULL,
  shop_id BIGINT NULL,
  customer_category_id BIGINT NULL,
  nickname VARCHAR(100) NULL,
  dob TIMESTAMP NOT NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0,
  user_id BIGINT NULL,
  last_txn_date TIMESTAMP NOT NULL
);

-- Create "identification_types" table
CREATE TABLE identification_types (
  identification_type_id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL DEFAULT '',
  code VARCHAR(100) NOT NULL DEFAULT '',
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);

-- Create "password_reset_tokens" table
CREATE TABLE password_reset_tokens (
  email VARCHAR(255) NOT NULL DEFAULT '',
  token VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL
);

-- Create "personal_access_token" table
CREATE TABLE personal_access_token (
  id BIGSERIAL PRIMARY KEY,
  tokenable_type VARCHAR(255) NOT NULL DEFAULT '',
  tokenable_id INTEGER NOT NULL DEFAULT 0,
  name VARCHAR(255) NOT NULL DEFAULT '',
  token VARCHAR(255) NOT NULL DEFAULT '',
  abilities VARCHAR(255) NOT NULL DEFAULT '',
  last_used_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL
);

-- Create "refresh_tokens" table
CREATE TABLE refresh_tokens (
  refresh_token_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL DEFAULT '',
  user_id BIGINT NOT NULL,
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

-- Create "roles" table
CREATE TABLE roles (
  role_id BIGSERIAL PRIMARY KEY,
  role VARCHAR(100) NOT NULL DEFAULT '',
  description VARCHAR(500) NOT NULL DEFAULT '',
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);

-- Create "shops" table
CREATE TABLE shops (
  shop_id BIGSERIAL PRIMARY KEY,
  shop_name VARCHAR(255) NOT NULL DEFAULT '',
  shop_description VARCHAR(255) NOT NULL DEFAULT '',
  shop_assistant_name VARCHAR(100) NOT NULL DEFAULT '',
  shop_assistant_number VARCHAR(100) NOT NULL DEFAULT '',
  image VARCHAR(100) NOT NULL DEFAULT '',
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);

-- Create "sign_up_d_t_o" table
CREATE TABLE sign_up_d_t_o (
  name VARCHAR(255) NOT NULL DEFAULT '',
  password VARCHAR(255) NOT NULL DEFAULT '',
  email VARCHAR(255) NOT NULL DEFAULT '',
  gender VARCHAR(255) NOT NULL DEFAULT '',
  dob VARCHAR(255) NOT NULL DEFAULT ''
);

-- Create "user_extra_details" table
CREATE TABLE user_extra_details (
  user_details_id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL DEFAULT 0,
  branch BIGINT NULL,
  shop_id BIGINT NULL,
  nickname VARCHAR(100) NULL,
  date_created TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);

-- Create "user_otps" table
CREATE TABLE user_otps (
  user_otp_id BIGSERIAL PRIMARY KEY,
  code VARCHAR(128) NOT NULL DEFAULT '',
  user_id BIGINT NOT NULL DEFAULT 0,
  status INTEGER NOT NULL DEFAULT 0,
  date_created TIMESTAMP NOT NULL,
  date_generated TIMESTAMP NOT NULL,
  expiry_date TIMESTAMP NOT NULL,
  date_modified TIMESTAMP NOT NULL,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 0
);

-- Create "user_tokens" table
CREATE TABLE user_tokens (
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

-- Create "username_d_t_o" table
CREATE TABLE username_d_t_o (
  username VARCHAR(255) NOT NULL DEFAULT ''
);
