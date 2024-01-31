package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Users_20230906_112843 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Users_20230906_112843{}
	m.Created = "20230906_112843"

	migration.Register("Users_20230906_112843", m)
}

// Run the migrations
func (m *Users_20230906_112843) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE users(`user_id` int(11) NOT NULL AUTO_INCREMENT,`user_type` int(11) DEFAULT NULL,`full_name` varchar(255) NOT NULL,`password` varchar(255) NOT NULL,`email` varchar(255) DEFAULT NULL,`phone_number` varchar(255) DEFAULT NULL,`gender` varchar(10) NOT NULL,`dob` datetime NOT NULL,`address` varchar(255) DEFAULT NULL,`id_type` varchar(5) DEFAULT NULL,`id_number` varchar(100) DEFAULT NULL,`active` int(11) DEFAULT 1,`is_verified` tinyint(1) DEFAULT 0,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,PRIMARY KEY (`user_id`))")
}

// Reverse the migrations
func (m *Users_20230906_112843) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `users`")
}
