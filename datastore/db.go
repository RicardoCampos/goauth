package datastore

import (
	"database/sql"
	// we mask the actual driver for now
	_ "github.com/lib/pq"
)

//NewDB creates a new database
func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
