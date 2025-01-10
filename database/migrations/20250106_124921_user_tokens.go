package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type UserTokens_20250106_124921 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserTokens_20250106_124921{}
	m.Created = "20250106_124921"

	migration.Register("UserTokens_20250106_124921", m)
}

// Run the migrations
func (m *UserTokens_20250106_124921) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE user_tokens(`user_token_id` int(11) NOT NULL AUTO_INCREMENT,`token` varchar(255) NOT NULL,`nonce` varchar(255) DEFAULT NULL,`expiry_date` datetime DEFAULT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,`active` int(11) DEFAULT 1,PRIMARY KEY (`user_token_id`))")
}

// Reverse the migrations
func (m *UserTokens_20250106_124921) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user_tokens`")
}
