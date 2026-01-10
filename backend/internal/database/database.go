package database

import (
	"database/sql"
	"log"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(database *sql.DB) {
	db = database
	log.Println("Database layer initialized")
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
