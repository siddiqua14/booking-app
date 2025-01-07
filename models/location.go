package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Location represents a location with associated hotels
type Location struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Hotels string `json:"hotels"` // Store hotels as a JSON string
}

// DB is the global database connection
var DB *sql.DB

// InitializeDatabase connects to the database and creates the locations table if it doesn't exist
func InitializeDatabase() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		"postgres",
		"5432",
		"postgres",
		"postgres",
		"booking_db",
		"disable",
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Create the locations table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS locations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		hotels JSONB NOT NULL
	)`
	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create locations table: %v", err)
	}

	log.Println("Database initialized successfully.")
	return nil
}

// AddLocation inserts a new location into the database
func AddLocation(name string, hotels string) error {
	query := `INSERT INTO locations (name, hotels) VALUES ($1, $2)`
	_, err := DB.Exec(query, name, hotels)
	return err
}
