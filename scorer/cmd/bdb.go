package main

import "database/sql"

type binDB struct {
	db   *sql.DB
	last string
	err  error
}

// Exec is a wrapper function which executes sql.DB.Exec and keep its error,
// when only the previous err is nil, which means there was no error before.
func (bdb *binDB) Exec(query string, args ...interface{}) {
	if bdb.err == nil {
		_, err := bdb.db.Exec(query, args...)
		bdb.last = query
		bdb.err = err
	}
}

// Prepare is a wrapper function which executes sql.DB.Prepare and keep its error,
// when only the previous err is nil, which means there was no error before.
func (bdb *binDB) Prepare(query string) *sql.Stmt {
	if bdb.err == nil {
		stmt, err := bdb.db.Prepare(query)
		bdb.last = query
		bdb.err = err
		return stmt
	}
	return nil
}
