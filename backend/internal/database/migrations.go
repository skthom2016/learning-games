package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// MigrationRecord tracks which migrations have been applied
type MigrationRecord struct {
	ID        int
	Filename  string
	AppliedAt string
}

// RunMigrations executes all pending migrations from the migrations directory
func RunMigrations(db *sql.DB, migrationsPath string) error {
	log.Println("Starting database migration check...")

	// Create schema_migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get all migration files
	migrationFiles, err := getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Sort migration files to ensure they run in order
	sort.Strings(migrationFiles)

	// Track if any migrations were applied
	appliedCount := 0

	// Execute pending migrations
	for _, filename := range migrationFiles {
		if _, exists := appliedMigrations[filename]; exists {
			continue // Skip already applied migrations
		}

		log.Printf("Applying migration: %s", filename)
		if err := applyMigration(db, filepath.Join(migrationsPath, filename), filename); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", filename, err)
		}
		appliedCount++
		log.Printf("Successfully applied migration: %s", filename)
	}

	if appliedCount == 0 {
		log.Println("No pending migrations found. Database is up to date.")
	} else {
		log.Printf("Successfully applied %d migration(s)", appliedCount)
	}

	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

// getAppliedMigrations returns a map of migration filenames that have been applied
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	query := "SELECT filename FROM schema_migrations ORDER BY id"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		applied[filename] = true
	}

	return applied, rows.Err()
}

// getMigrationFiles returns a sorted list of SQL migration files
func getMigrationFiles(migrationsPath string) ([]string, error) {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	return migrationFiles, nil
}

// applyMigration executes a single migration file within a transaction
func applyMigration(db *sql.DB, filepath string, filename string) error {
	// Read migration file
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Rollback on error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Execute migration SQL
	if _, err = tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as applied
	insertQuery := "INSERT INTO schema_migrations (filename) VALUES ($1)"
	if _, err = tx.Exec(insertQuery, filename); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
