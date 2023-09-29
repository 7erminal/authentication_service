package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type UserTypes_20230906_101749 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserTypes_20230906_101749{}
	m.Created = "20230906_101749"

	migration.Register("UserTypes_20230906_101749", m)
}

// Run the migrations
func (m *UserTypes_20230906_101749) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE user_types(`user_type_id` int(11) NOT NULL AUTO_INCREMENT,`user_type_name` varchar(255) NOT NULL,`user_type_description` varchar(255) NOT NULL,`active` int(11) DEFAULT 1,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,PRIMARY KEY (`user_type_id`))")
}

// Reverse the migrations
func (m *UserTypes_20230906_101749) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user_types`")
}
