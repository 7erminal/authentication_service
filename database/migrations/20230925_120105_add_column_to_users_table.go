package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddColumnToUsersTable_20230925_120105 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddColumnToUsersTable_20230925_120105{}
	m.Created = "20230925_120105"

	migration.Register("AddColumnToUsersTable_20230925_120105", m)
}

// Run the migrations
func (m *AddColumnToUsersTable_20230925_120105) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE users ADD COLUMN `username` varchar(40) DEFAULT NULL AFTER `full_name`")
}

// Reverse the migrations
func (m *AddColumnToUsersTable_20230925_120105) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
