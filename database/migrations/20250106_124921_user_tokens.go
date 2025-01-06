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
	m.SQL("CREATE TABLE user_tokens(`id` int(11) NOT NULL AUTO_INCREMENT,`token` varchar(255) NOT NULL,`nonce` varchar(255) NOT NULL,`expiry_date` datetime NOT NULL,`date_created` datetime NOT NULL,`date_modified` datetime NOT NULL,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,`active` int(11) DEFAULT NULL,PRIMARY KEY (`id`))")
}

// Reverse the migrations
func (m *UserTokens_20250106_124921) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user_tokens`")
}
