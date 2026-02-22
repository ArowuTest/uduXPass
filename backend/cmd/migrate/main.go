package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	SQL         string
	FilePath    string
}

func main() {
	var (
		dbURL       = flag.String("db-url", "", "Database URL (required)")
		migrationsDir = flag.String("migrations-dir", "./migrations", "Migrations directory")
		action      = flag.String("action", "up", "Migration action: up, down, status")
		steps       = flag.Int("steps", 0, "Number of migration steps (0 = all)")
	)
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("Database URL is required. Use -db-url flag or set DATABASE_URL environment variable")
	}

	// Connect to database
	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Load migrations
	migrations, err := loadMigrations(*migrationsDir)
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	// Execute action
	switch *action {
	case "up":
		if err := migrateUp(db, migrations, *steps); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
	case "down":
		if err := migrateDown(db, migrations, *steps); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
	case "status":
		if err := showStatus(db, migrations); err != nil {
			log.Fatalf("Failed to show status: %v", err)
		}
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

// createMigrationsTable creates the migrations tracking table
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`
	
	_, err := db.Exec(query)
	return err
}

// loadMigrations loads all migration files from the directory
func loadMigrations(dir string) ([]Migration, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse filename: 001_initial_schema.sql
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		description := strings.TrimSuffix(parts[1], ".sql")
		description = strings.ReplaceAll(description, "_", " ")

		filePath := filepath.Join(dir, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filePath, err)
		}

		migrations = append(migrations, Migration{
			Version:     version,
			Description: description,
			SQL:         string(content),
			FilePath:    filePath,
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns a map of applied migration versions
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	query := "SELECT version FROM schema_migrations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// migrateUp applies pending migrations
func migrateUp(db *sql.DB, migrations []Migration, steps int) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	var pending []Migration
	for _, migration := range migrations {
		if !applied[migration.Version] {
			pending = append(pending, migration)
		}
	}

	if len(pending) == 0 {
		fmt.Println("No pending migrations")
		return nil
	}

	// Limit steps if specified
	if steps > 0 && steps < len(pending) {
		pending = pending[:steps]
	}

	fmt.Printf("Applying %d migration(s):\n", len(pending))

	for _, migration := range pending {
		fmt.Printf("  %s: %s\n", migration.Version, migration.Description)

		// Execute migration SQL
		if _, err := db.Exec(migration.SQL); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}

		// Record migration as applied
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		fmt.Printf("  ✓ Applied %s\n", migration.Version)
	}

	fmt.Println("Migration completed successfully")
	return nil
}

// migrateDown rolls back migrations
func migrateDown(db *sql.DB, migrations []Migration, steps int) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find applied migrations in reverse order
	var toRollback []Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		if applied[migration.Version] {
			toRollback = append(toRollback, migration)
		}
	}

	if len(toRollback) == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Limit steps if specified
	if steps > 0 && steps < len(toRollback) {
		toRollback = toRollback[:steps]
	}

	fmt.Printf("Rolling back %d migration(s):\n", len(toRollback))

	for _, migration := range toRollback {
		fmt.Printf("  %s: %s\n", migration.Version, migration.Description)

		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for rollback %s: %w", migration.Version, err)
		}

		// Note: This is a simple rollback that just removes the migration record
		// In a production system, you'd want to have separate down migration files
		if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %s: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit rollback %s: %w", migration.Version, err)
		}

		fmt.Printf("  ✓ Rolled back %s\n", migration.Version)
	}

	fmt.Println("Rollback completed successfully")
	return nil
}

// showStatus displays the current migration status
func showStatus(db *sql.DB, migrations []Migration) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, migration := range migrations {
		status := "PENDING"
		if applied[migration.Version] {
			status = "APPLIED"
		}
		fmt.Printf("  %s: %s [%s]\n", migration.Version, migration.Description, status)
	}

	appliedCount := len(applied)
	totalCount := len(migrations)
	pendingCount := totalCount - appliedCount

	fmt.Printf("\nSummary: %d applied, %d pending, %d total\n", appliedCount, pendingCount, totalCount)

	return nil
}

