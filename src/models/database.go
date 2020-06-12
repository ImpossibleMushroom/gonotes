package models

import (
	"database/sql"
	"fmt"

	// Postgres driver for database/sql
	_ "github.com/lib/pq"
)

// Database details
//	In a production environment, pass these through command line args or environment variables
const (
	host     = "localhost"
	port     = 5432
	user     = "gonotes"
	password = "root"
	dbname   = "GoNotes"
)

// Database global variable, not exported
var db *sql.DB

// NoRows - Dirty way to get access to sql.ErrNoRows without importing database/sql in other models
var NoRows = sql.ErrNoRows

// ConnectToDatabase ...
func ConnectToDatabase() error {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	database, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return err
	}

	db = database

	return nil
}
