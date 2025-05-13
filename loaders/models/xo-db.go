package models

import "database/sql"

// IDB is the common interface for database operations that can be used with
// types from schema 'public'.
//
// This should work with database/sql.DB and database/sql.Tx.
type IDB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}
