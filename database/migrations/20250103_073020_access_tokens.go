package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AccessTokens_20250103_073020 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AccessTokens_20250103_073020{}
	m.Created = "20250103_073020"

	migration.Register("AccessTokens_20250103_073020", m)
}

// Run the migrations
func (m *AccessTokens_20250103_073020) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE access_tokens(`access_token_id` int(11) NOT NULL AUTO_INCREMENT,`token` varchar(255) NOT NULL,`user_id` int(11) NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`expires_at` datetime NOT NULL,`revoked` tinyint(1) DEFAULT 0,`ip_address` varchar(80) DEFAULT NULL,`last_used_at` datetime DEFAULT NULL,PRIMARY KEY (`access_token_id`), FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE)")
}

// Reverse the migrations
func (m *AccessTokens_20250103_073020) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `access_tokens`")
}
