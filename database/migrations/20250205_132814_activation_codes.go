package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type ActivationCodes_20250205_132814 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ActivationCodes_20250205_132814{}
	m.Created = "20250205_132814"

	migration.Register("ActivationCodes_20250205_132814", m)
}

// Run the migrations
func (m *ActivationCodes_20250205_132814) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE activation_codes(`activation_code_id` int(11) NOT NULL AUTO_INCREMENT,`code` varchar(80) NOT NULL,`number` varchar(80) NOT NULL,`expiry_date` datetime NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,`active` int(11) DEFAULT 1,PRIMARY KEY (`activation_code_id`))")
}

// Reverse the migrations
func (m *ActivationCodes_20250205_132814) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `activation_codes`")
}
