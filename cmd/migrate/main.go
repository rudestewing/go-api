package main

import (
	"fmt"
	"go-api/config"
	"go-api/internal/migration"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	// Initialize config
	config.InitConfig()

	command := os.Args[1]
	switch command {
	case "migrate":
		runMigrations()
	case "rollback":
		rollbackMigrations()
	case "status":
		showStatus()
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Migration name is required. Usage: go run cmd/migrate/main.go create \"migration_name\"")
		}
		createMigration(os.Args[2])
	case "help", "--help", "-h":
		showUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showUsage()
	}
}

func runMigrations() {
	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func rollbackMigrations() {
	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.RollbackLastBatch(); err != nil {
		log.Fatalf("Failed to rollback migrations: %v", err)
	}
}

func showStatus() {
	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.GetMigrationStatus(); err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}
}

func createMigration(name string) {
	if name == "" {
		log.Fatal("Migration name is required. Use -name flag")
	}

	// Create migrations directory if it doesn't exist
	migrationsDir := "migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Generate timestamp
	timestamp := time.Now().Format("20060102150405")

	// Clean migration name (replace spaces with underscores, remove special chars)
	cleanName := strings.ReplaceAll(name, " ", "_")
	cleanName = strings.ToLower(cleanName)

	// Create filename
	filename := fmt.Sprintf("%s_%s.sql", timestamp, cleanName)
	filepath := filepath.Join(migrationsDir, filename)

	// Create migration file template
	template := `-- +migrate Up
-- Write your UP migration here
-- Example:
-- CREATE TABLE example (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- +migrate Down
-- Write your DOWN migration here
-- Example:
-- DROP TABLE IF EXISTS example;
`

	if err := os.WriteFile(filepath, []byte(template), 0644); err != nil {
		log.Fatalf("Failed to create migration file: %v", err)
	}

	fmt.Printf("Migration file created: %s\n", filepath)
	fmt.Println("Please edit the file to add your migration SQL")
}

func showUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println("=======================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go create \"migration_name\"")
	fmt.Println("  go run cmd/migrate/main.go migrate")
	fmt.Println("  go run cmd/migrate/main.go rollback")
	fmt.Println("  go run cmd/migrate/main.go status")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create <name> - Create a new migration file")
	fmt.Println("  migrate       - Run all pending migrations")
	fmt.Println("  rollback      - Rollback the last batch of migrations")
	fmt.Println("  status        - Show migration status")
	fmt.Println("  help          - Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go create \"create_users_table\"")
	fmt.Println("  go run cmd/migrate/main.go migrate")
	fmt.Println("  go run cmd/migrate/main.go rollback")
	fmt.Println("  go run cmd/migrate/main.go status")
}
