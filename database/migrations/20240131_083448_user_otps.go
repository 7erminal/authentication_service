package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type UserOtps_20240131_083448 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserOtps_20240131_083448{}
	m.Created = "20240131_083448"

	migration.Register("UserOtps_20240131_083448", m)
}

// Run the migrations
func (m *UserOtps_20240131_083448) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE user_otps(`user_otp_id` int(11) NOT NULL AUTO_INCREMENT,`user_id` int(11) NOT NULL,`one_time_pin` int(6) NOT NULL,`active` int(11) DEFAULT 0,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,PRIMARY KEY (`user_otp_id`), FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *UserOtps_20240131_083448) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user_otps`")
}
