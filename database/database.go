package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cat-api/config"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(cfg *config.Config) error {
	// Create database directory if it doesn't exist
	dbDir := filepath.Dir(cfg.DBPath)
	if dbDir != "." && dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %v", err)
		}
	}

	var err error
	DB, err = sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// SQLite from modernc.org/sqlite automatically supports foreign keys
	// No need for PRAGMA foreign_keys = ON

	// Enable WAL mode for better concurrency
	if _, err := DB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		log.Printf("Warning: Failed to enable WAL mode: %v", err)
	}

	// Set busy timeout
	if _, err := DB.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		log.Printf("Warning: Failed to set busy timeout: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Connected to SQLite database: %s", cfg.DBPath)
	return runMigrations()
}

// Остальной код без изменений...
func runMigrations() error {
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
        CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `

	if _, err := DB.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Check if initial migration has been applied
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = '001_init'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migrations: %v", err)
	}

	if count == 0 {
		log.Println("Applying initial migration...")

		// Read migration file
		migrationSQL, err := os.ReadFile("database/migrations/001_init.sql")
		if err != nil {
			return fmt.Errorf("failed to read migration file: %v", err)
		}

		// Start transaction
		tx, err := DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}
		defer tx.Rollback()

		// Execute migration
		if _, err := tx.Exec(string(migrationSQL)); err != nil {
			return fmt.Errorf("failed to execute migration: %v", err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO migrations (name) VALUES ('001_init')"); err != nil {
			return fmt.Errorf("failed to record migration: %v", err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration: %v", err)
		}

		log.Println("Initial migration applied successfully")
	} else {
		log.Println("Migrations are up to date")
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}
