//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Users = newUsersTable("public", "users", "")

type usersTable struct {
	postgres.Table

	//Columns
	UserID   postgres.ColumnString
	Username postgres.ColumnString
	Password postgres.ColumnString
	UserRole postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UsersTable struct {
	usersTable

	EXCLUDED usersTable
}

// AS creates new UsersTable with assigned alias
func (a UsersTable) AS(alias string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UsersTable with assigned schema name
func (a UsersTable) FromSchema(schemaName string) *UsersTable {
	return newUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UsersTable with assigned table prefix
func (a UsersTable) WithPrefix(prefix string) *UsersTable {
	return newUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UsersTable with assigned table suffix
func (a UsersTable) WithSuffix(suffix string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUsersTable(schemaName, tableName, alias string) *UsersTable {
	return &UsersTable{
		usersTable: newUsersTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newUsersTableImpl("", "excluded", ""),
	}
}

func newUsersTableImpl(schemaName, tableName, alias string) usersTable {
	var (
		UserIDColumn   = postgres.StringColumn("user_id")
		UsernameColumn = postgres.StringColumn("username")
		PasswordColumn = postgres.StringColumn("password")
		UserRoleColumn = postgres.StringColumn("user_role")
		allColumns     = postgres.ColumnList{UserIDColumn, UsernameColumn, PasswordColumn, UserRoleColumn}
		mutableColumns = postgres.ColumnList{UsernameColumn, PasswordColumn, UserRoleColumn}
	)

	return usersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:   UserIDColumn,
		Username: UsernameColumn,
		Password: PasswordColumn,
		UserRole: UserRoleColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
