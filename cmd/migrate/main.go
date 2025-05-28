package main

import (
	"fmt"
	"go-api/config"
	"go-api/database/migration"
	"log"
	"os"
	"strconv"
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
	case "migrate", "up":
		runMigrations()
	case "rollback", "down":
		rollbackMigrations()
	case "status":
		showStatus()
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Migration name is required. Usage: go run cmd/migrate/main.go create \"migration_name\"")
		}
		createMigration(os.Args[2])
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Version is required. Usage: go run cmd/migrate/main.go force <version>")
		}
		forceVersion(os.Args[2])
	case "drop":
		dropDatabase()
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

	// Check if user wants to rollback to a specific version
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "last", "1":
			if err := manager.RollbackLastMigration(); err != nil {
				log.Fatalf("Failed to rollback last migration: %v", err)
			}
		default:
			log.Fatal("Usage: go run cmd/migrate/main.go rollback [last|1]")
		}
	} else {
		if err := manager.RollbackLastMigration(); err != nil {
			log.Fatalf("Failed to rollback last migration: %v", err)
		}
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
	if err := migration.CreateMigration(name); err != nil {
		log.Fatalf("Failed to create migration: %v", err)
	}
}

func forceVersion(versionStr string) {
	// Import strconv for string to int conversion
	version := 0
	if versionStr != "0" {
		var err error
		version, err = strconv.Atoi(versionStr)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
	}

	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.Force(version); err != nil {
		log.Fatalf("Failed to force version: %v", err)
	}
}

func dropDatabase() {
	fmt.Print("⚠️  This will drop all tables in the database. Are you sure? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" && response != "yes" && response != "YES" {
		fmt.Println("Operation cancelled")
		return
	}

	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.Drop(); err != nil {
		log.Fatalf("Failed to drop database: %v", err)
	}
}

func showUsage() {
	fmt.Println("Database Migration Tool (golang-migrate)")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create <name>  - Create a new migration file pair (.up.sql and .down.sql)")
	fmt.Println("  migrate/up     - Run all pending migrations")
	fmt.Println("  rollback/down  - Rollback the last migration")
	fmt.Println("  status         - Show current migration status")
	fmt.Println("  force <version>- Force set the migration version (use with caution)")
	fmt.Println("  drop           - Drop all tables (DANGEROUS - requires confirmation)")
	fmt.Println("  help           - Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go create \"create_users_table\"")
	fmt.Println("  go run cmd/migrate/main.go migrate")
	fmt.Println("  go run cmd/migrate/main.go rollback")
	fmt.Println("  go run cmd/migrate/main.go status")
	fmt.Println("  go run cmd/migrate/main.go force 3")
	fmt.Println("  go run cmd/migrate/main.go drop")
	fmt.Println()
	fmt.Println("Note: Migration files are now in separate .up.sql and .down.sql files")
	fmt.Println("      This follows the golang-migrate library standard")
}
