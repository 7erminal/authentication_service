-- Create "customer_refresh_tokens" table
CREATE TABLE `customer_refresh_tokens` (
  `refresh_token_id` bigint NOT NULL AUTO_INCREMENT,
  `token` varchar(255) NOT NULL DEFAULT "",
  `customer_id` bigint NOT NULL,
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
