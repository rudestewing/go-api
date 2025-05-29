package migration

import (
	"database/sql"
	"fmt"
	"go-api/config"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type MigrationManager struct {
	migrate *migrate.Migrate
	db      *sql.DB
}

func NewMigrationManager() (*MigrationManager, error) {
	cfg := config.Get()
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// Open database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get absolute path to migrations directory
	migrationsPath, err := filepath.Abs("database/migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path to migrations: %w", err)
	}

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &MigrationManager{
		migrate: m,
		db:      db,
	}, nil
}

func (m *MigrationManager) RunMigrations() error {
	err := m.migrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No pending migrations to run")
	} else {
		log.Println("All migrations completed successfully")
	}

	return nil
}

func (m *MigrationManager) RollbackLastMigration() error {
	err := m.migrate.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No migrations to rollback")
	} else {
		log.Println("Migration rolled back successfully")
	}

	return nil
}

func (m *MigrationManager) RollbackToVersion(version uint) error {
	err := m.migrate.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	if err == migrate.ErrNoChange {
		log.Printf("Already at version %d", version)
	} else {
		log.Printf("Migrated to version %d successfully", version)
	}

	return nil
}

func (m *MigrationManager) GetMigrationStatus() error {
	version, dirty, err := m.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	if err == migrate.ErrNilVersion {
		fmt.Println("No migrations have been applied")
	} else {
		fmt.Printf("Current version: %d\n", version)
		if dirty {
			fmt.Println("State: DIRTY (migration failed or was interrupted)")
			fmt.Println("⚠️  Please fix the dirty state before running new migrations")
		} else {
			fmt.Println("State: CLEAN")
		}
	}

	// List available migrations
	if err := m.listAvailableMigrations(); err != nil {
		log.Printf("Warning: Failed to list available migrations: %v", err)
	}

	return nil
}

func (m *MigrationManager) listAvailableMigrations() error {
	migrationsDir := "database/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	upMigrations := make(map[string]bool)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationName := strings.TrimSuffix(file.Name(), ".up.sql")
			upMigrations[migrationName] = true
		}
	}

	if len(upMigrations) > 0 {
		fmt.Println("\nAvailable migrations:")
		for migration := range upMigrations {
			fmt.Printf("  📄 %s\n", migration)
		}
	}

	return nil
}

func (m *MigrationManager) Force(version int) error {
	err := m.migrate.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force version %d: %w", version, err)
	}

	log.Printf("Forced migration to version %d", version)
	return nil
}

func (m *MigrationManager) Drop() error {
	err := m.migrate.Drop()
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	log.Println("Database dropped successfully")
	return nil
}

// Fresh drops all tables and re-runs all migrations from the beginning
func (m *MigrationManager) Fresh() error {
	log.Println("Starting fresh migration...")

	// First, drop all tables
	if err := m.migrate.Drop(); err != nil {
		return fmt.Errorf("failed to drop database during fresh: %w", err)
	}
	log.Println("All tables dropped")

	// After dropping all tables, we need to create a new migration instance
	// because the schema_migrations table was also dropped
	freshManager, err := NewMigrationManager()
	if err != nil {
		return fmt.Errorf("failed to create fresh migration manager: %w", err)
	}
	defer freshManager.Close()

	// Run all migrations from the beginning with the fresh manager
	if err := freshManager.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations during fresh: %w", err)
	}

	log.Println("Fresh migration completed successfully - all migrations re-run from beginning")
	return nil
}

// Purge rolls back all executed migrations to version 0
func (m *MigrationManager) Purge() error {
	log.Println("Starting migration purge...")

	// Get current version
	version, dirty, err := m.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations to purge - database is already at version 0")
		return nil
	}

	if dirty {
		return fmt.Errorf("cannot purge migrations: database is in dirty state (version %d). Please fix the dirty state first using 'force' command", version)
	}

	initialVersion := version
	log.Printf("Starting purge from version %d", initialVersion)

	// Roll back all migrations one by one until we reach version 0
	for {
		currentVersion, _, err := m.migrate.Version()
		if err == migrate.ErrNilVersion {
			// We've successfully rolled back all migrations
			break
		}
		if err != nil {
			return fmt.Errorf("failed to get current version during purge: %w", err)
		}

		// Perform one step rollback
		log.Printf("Rolling back from version %d...", currentVersion)
		if err := m.migrate.Steps(-1); err != nil {
			if err == migrate.ErrNoChange {
				// No more migrations to rollback
				break
			}
			return fmt.Errorf("failed to rollback step during purge: %w", err)
		}
	}

	log.Printf("Migration purge completed - rolled back from version %d to version 0", initialVersion)
	return nil
}

func (m *MigrationManager) Close() error {
	if m.migrate != nil {
		if sourceErr, dbErr := m.migrate.Close(); sourceErr != nil || dbErr != nil {
			return fmt.Errorf("failed to close migrate instance: source=%v, db=%v", sourceErr, dbErr)
		}
	}

	if m.db != nil {
		return m.db.Close()
	}

	return nil
}

// CreateMigration creates a new migration file pair (.up.sql and .down.sql)
func CreateMigration(name string) error {
	if name == "" {
		return fmt.Errorf("migration name is required")
	}

	// Create migrations directory if it doesn't exist
	migrationsDir := "database/migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Generate timestamp
	timestamp := time.Now().Format("20060102150405")

	// Clean migration name (replace spaces with underscores, remove special chars)
	cleanName := strings.ReplaceAll(name, " ", "_")
	cleanName = strings.ToLower(cleanName)

	// Create filenames
	upFilename := fmt.Sprintf("%s_%s.up.sql", timestamp, cleanName)
	downFilename := fmt.Sprintf("%s_%s.down.sql", timestamp, cleanName)

	upFilepath := filepath.Join(migrationsDir, upFilename)
	downFilepath := filepath.Join(migrationsDir, downFilename)

	// Create UP migration file template
	upTemplate := `-- Write your UP migration here
-- Example:
-- CREATE TABLE example (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
`

	// Create DOWN migration file template
	downTemplate := `-- Write your DOWN migration here
-- Example:
-- DROP TABLE IF EXISTS example;
`

	// Write UP migration file
	if err := os.WriteFile(upFilepath, []byte(upTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create UP migration file: %w", err)
	}

	// Write DOWN migration file
	if err := os.WriteFile(downFilepath, []byte(downTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create DOWN migration file: %w", err)
	}

	fmt.Printf("Migration files created:\n")
	fmt.Printf("  📝 %s\n", upFilepath)
	fmt.Printf("  📝 %s\n", downFilepath)
	fmt.Println("Please edit the files to add your migration SQL")

	return nil
}
