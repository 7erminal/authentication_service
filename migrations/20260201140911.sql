-- Create "access_tokens" table
CREATE TABLE `access_tokens` (
  `access_token_id` bigint NOT NULL AUTO_INCREMENT,
  `token` varchar(255) NOT NULL DEFAULT "",
  `user_id` bigint NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `expires_at` datetime NOT NULL,
  `revoked` bool NOT NULL DEFAULT 0,
  `i_p_address` varchar(80) NOT NULL DEFAULT "",
  `last_used_at` datetime NOT NULL,
  PRIMARY KEY (`access_token_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "activation_codes" table
CREATE TABLE `activation_codes` (
  `activation_code_id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(80) NOT NULL DEFAULT "",
  `number` varchar(80) NOT NULL DEFAULT "",
  `expiry_date` datetime NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`activation_code_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "auth_users" table
CREATE TABLE `auth_users` (
  `user_id` bigint NOT NULL AUTO_INCREMENT,
  `user_type` int NOT NULL DEFAULT 0,
  `user_details_id` bigint NULL,
  `image_path` varchar(200) NULL,
  `full_name` varchar(255) NOT NULL DEFAULT "",
  `username` varchar(255) NOT NULL DEFAULT "",
  `password` varchar(255) NOT NULL DEFAULT "",
  `email` varchar(255) NOT NULL DEFAULT "",
  `phone_number` varchar(255) NOT NULL DEFAULT "",
  `gender` varchar(10) NOT NULL DEFAULT "",
  `dob` datetime NOT NULL,
  `address` varchar(255) NOT NULL DEFAULT "",
  `id_type` varchar(5) NOT NULL DEFAULT "",
  `id_number` varchar(100) NOT NULL DEFAULT "",
  `marital_status` varchar(255) NOT NULL DEFAULT "",
  `role` bigint NULL,
  `active` int NOT NULL DEFAULT 0,
  `is_verified` bool NOT NULL DEFAULT 0,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "authentication_d_t_o" table
CREATE TABLE `authentication_d_t_o` (
  `username` varchar(255) NOT NULL DEFAULT "",
  `password` varchar(255) NOT NULL DEFAULT ""
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "branches" table
CREATE TABLE `branches` (
  `branch_id` bigint NOT NULL AUTO_INCREMENT,
  `branch` varchar(80) NOT NULL DEFAULT "",
  `country_id` bigint NOT NULL,
  `location` varchar(255) NOT NULL DEFAULT "",
  `phone_number` varchar(255) NOT NULL DEFAULT "",
  `active` int NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NULL,
  `modified_by` int NULL,
  PRIMARY KEY (`branch_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "countries" table
CREATE TABLE `countries` (
  `country_id` bigint NOT NULL AUTO_INCREMENT,
  `country` varchar(255) NOT NULL DEFAULT "",
  `description` varchar(500) NOT NULL DEFAULT "",
  `country_code` varchar(20) NOT NULL DEFAULT "",
  `default_currency` bigint NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`country_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "currencies" table
CREATE TABLE `currencies` (
  `currency_id` bigint NOT NULL AUTO_INCREMENT,
  `symbol` varchar(20) NOT NULL DEFAULT "",
  `currency` varchar(50) NOT NULL DEFAULT "",
  `active` int NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NULL,
  `modified_by` int NULL,
  PRIMARY KEY (`currency_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "customer_access_tokens" table
CREATE TABLE `customer_access_tokens` (
  `customer_access_token_id` bigint NOT NULL AUTO_INCREMENT,
  `token` varchar(255) NOT NULL DEFAULT "",
  `customer_id` bigint NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `expires_at` datetime NOT NULL,
  `revoked` bool NOT NULL DEFAULT 0,
  `ip_address` varchar(80) NOT NULL DEFAULT "",
  `last_used_at` datetime NOT NULL,
  PRIMARY KEY (`customer_access_token_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "customer_categories" table
CREATE TABLE `customer_categories` (
  `customer_category_id` bigint NOT NULL AUTO_INCREMENT,
  `category` varchar(100) NOT NULL DEFAULT "",
  `description` varchar(255) NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`customer_category_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "customer_credentials" table
CREATE TABLE `customer_credentials` (
  `customer_credential_id` bigint NOT NULL AUTO_INCREMENT,
  `customer_id` bigint NOT NULL,
  `username` varchar(255) NOT NULL DEFAULT "",
  `password` varchar(255) NOT NULL DEFAULT "",
  `pin` varchar(10) NOT NULL DEFAULT "",
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`customer_credential_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "customer_emergency_contacts" table
CREATE TABLE `customer_emergency_contacts` (
  `customer_emergency_contact_id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(120) NOT NULL DEFAULT "",
  `contact` varchar(50) NOT NULL DEFAULT "",
  `customer_id` bigint NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`customer_emergency_contact_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "customer_guarantors" table
CREATE TABLE `customer_guarantors` (
  `customer_guarantor_id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(120) NOT NULL DEFAULT "",
  `contact` varchar(50) NOT NULL DEFAULT "",
  `customer_id` bigint NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`customer_guarantor_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "customers" table
CREATE TABLE `customers` (
  `customer_id` bigint NOT NULL AUTO_INCREMENT,
  `customer_number` varchar(255) NOT NULL DEFAULT "",
  `full_name` varchar(255) NOT NULL DEFAULT "",
  `image_path` varchar(255) NOT NULL DEFAULT "",
  `email` varchar(255) NULL,
  `phone_number` varchar(255) NULL,
  `location` varchar(255) NULL,
  `identification_type_id` bigint NULL,
  `identification_number` varchar(255) NULL,
  `branch` bigint NULL,
  `shop_id` bigint NULL,
  `customer_category_id` bigint NULL,
  `nickname` varchar(100) NULL,
  `dob` datetime NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  `user_id` bigint NULL,
  `last_txn_date` datetime NOT NULL,
  PRIMARY KEY (`customer_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "identification_types" table
CREATE TABLE `identification_types` (
  `identification_type_id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT "",
  `code` varchar(100) NOT NULL DEFAULT "",
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`identification_type_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "password_reset_tokens" table
CREATE TABLE `password_reset_tokens` (
  `email` varchar(255) NOT NULL DEFAULT "",
  `token` varchar(255) NOT NULL DEFAULT "",
  `created_at` datetime NOT NULL
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "personal_access_token" table
CREATE TABLE `personal_access_token` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `tokenable_type` varchar(255) NOT NULL DEFAULT "",
  `tokenable_id` int NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL DEFAULT "",
  `token` varchar(255) NOT NULL DEFAULT "",
  `abilities` varchar(255) NOT NULL DEFAULT "",
  `last_used_at` datetime NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "refresh_tokens" table
CREATE TABLE `refresh_tokens` (
  `refresh_token_id` bigint NOT NULL AUTO_INCREMENT,
  `token` varchar(255) NOT NULL DEFAULT "",
  `user_id` bigint NOT NULL,
  `access_token_id` bigint NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `expires_at` datetime NOT NULL,
  `revoked` bool NOT NULL DEFAULT 0,
  `i_p_address` varchar(45) NULL,
  `user_agent` varchar(255) NULL,
  `last_used_at` datetime NULL,
  PRIMARY KEY (`refresh_token_id`),
  UNIQUE INDEX `token` (`token`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "roles" table
CREATE TABLE `roles` (
  `role_id` bigint NOT NULL AUTO_INCREMENT,
  `role` varchar(100) NOT NULL DEFAULT "",
  `description` varchar(500) NOT NULL DEFAULT "",
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`role_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "shops" table
CREATE TABLE `shops` (
  `shop_id` bigint NOT NULL AUTO_INCREMENT,
  `shop_name` varchar(255) NOT NULL DEFAULT "",
  `shop_description` varchar(255) NOT NULL DEFAULT "",
  `shop_assistant_name` varchar(100) NOT NULL DEFAULT "",
  `shop_assistant_number` varchar(100) NOT NULL DEFAULT "",
  `image` varchar(100) NOT NULL DEFAULT "",
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`shop_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "sign_up_d_t_o" table
CREATE TABLE `sign_up_d_t_o` (
  `name` varchar(255) NOT NULL DEFAULT "",
  `password` varchar(255) NOT NULL DEFAULT "",
  `email` varchar(255) NOT NULL DEFAULT "",
  `gender` varchar(255) NOT NULL DEFAULT "",
  `dob` varchar(255) NOT NULL DEFAULT ""
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "user_extra_details" table
CREATE TABLE `user_extra_details` (
  `user_details_id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL DEFAULT 0,
  `branch` bigint NULL,
  `shop_id` bigint NULL,
  `nickname` varchar(100) NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_details_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "user_otps" table
CREATE TABLE `user_otps` (
  `user_otp_id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(128) NOT NULL DEFAULT "",
  `user_id` bigint NOT NULL DEFAULT 0,
  `status` int NOT NULL DEFAULT 0,
  `date_created` datetime NOT NULL,
  `date_generated` datetime NOT NULL,
  `expiry_date` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_otp_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "user_tokens" table
CREATE TABLE `user_tokens` (
  `user_token_id` bigint NOT NULL AUTO_INCREMENT,
  `token` varchar(255) NOT NULL DEFAULT "",
  `nonce` varchar(255) NOT NULL DEFAULT "",
  `expiry_date` datetime NOT NULL,
  `date_created` datetime NOT NULL,
  `date_modified` datetime NOT NULL,
  `created_by` int NOT NULL DEFAULT 0,
  `modified_by` int NOT NULL DEFAULT 0,
  `active` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_token_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "username_d_t_o" table
CREATE TABLE `username_d_t_o` (
  `username` varchar(255) NOT NULL DEFAULT ""
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
