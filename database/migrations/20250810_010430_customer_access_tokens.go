package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CustomerAccessTokens_20250810_010430 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CustomerAccessTokens_20250810_010430{}
	m.Created = "20250810_010430"

	migration.Register("CustomerAccessTokens_20250810_010430", m)
}

// Run the migrations
func (m *CustomerAccessTokens_20250810_010430) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE customer_access_tokens(`customer_access_token_id` int(11) NOT NULL AUTO_INCREMENT,`token` varchar(255) NOT NULL,`customer_id` int(11) NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`expires_at` datetime DEFAULT NULL,`revoked` tinyint(1) DEFAULT 0,`ip_address` varchar(80) DEFAULT NULL,`last_used_at` datetime NOT NULL,PRIMARY KEY (`customer_access_token_id`), FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON UPDATE CASCADE ON DELETE CASCADE)")
}

// Reverse the migrations
func (m *CustomerAccessTokens_20250810_010430) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `customer_access_tokens`")
}
