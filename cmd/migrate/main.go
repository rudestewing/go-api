package main

import (
	"bufio"
	"fmt"
	"go-api/config"
	"go-api/database/migration"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	// Initialize config
	config.InitConfig()

	command := os.Args[1]

	// Validate command before processing
	if !isValidCommand(command) {
		fmt.Printf("Unknown command: %s\n\n", command)
		showUsage()
		return
	}

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
		if err := validateMigrationName(os.Args[2]); err != nil {
			log.Fatalf("Invalid migration name: %v", err)
		}
		createMigration(os.Args[2])
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Version is required. Usage: go run cmd/migrate/main.go force <version>")
		}
		version, err := validateAndParseVersion(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		forceVersion(version)
	case "drop":
		if err := confirmDangerousOperation("drop all database tables"); err != nil {
			log.Fatalf("Operation cancelled: %v", err)
		}
		dropDatabase()
	case "fresh":
		if err := confirmDangerousOperation("drop and recreate all database tables"); err != nil {
			log.Fatalf("Operation cancelled: %v", err)
		}
		runFreshMigrations()
	case "purge":
		if err := confirmDangerousOperation("purge all migrations"); err != nil {
			log.Fatalf("Operation cancelled: %v", err)
		}
		purgeMigrations()
	case "help", "--help", "-h":
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

func forceVersion(version int) {
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
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

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

func runFreshMigrations() {
	fmt.Print("⚠️  This will drop all tables and re-run all migrations from the beginning. Are you sure? (y/N): ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	if response != "y" && response != "Y" && response != "yes" && response != "YES" {
		fmt.Println("Operation cancelled")
		return
	}

	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.Fresh(); err != nil {
		log.Fatalf("Failed to run fresh migrations: %v", err)
	}
}

func purgeMigrations() {
	fmt.Print("⚠️  This will rollback all executed migrations to version 0. Are you sure? (y/N): ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	if response != "y" && response != "Y" && response != "yes" && response != "YES" {
		fmt.Println("Operation cancelled")
		return
	}

	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.Purge(); err != nil {
		log.Fatalf("Failed to purge migrations: %v", err)
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
	fmt.Println("  create <n>     - Create a new migration file pair (.up.sql and .down.sql)")
	fmt.Println("  migrate/up     - Run all pending migrations")
	fmt.Println("  rollback/down  - Rollback the last migration")
	fmt.Println("  fresh          - Drop all tables and re-run all migrations from beginning")
	fmt.Println("  purge          - Rollback all executed migrations to version 0")
	fmt.Println("  status         - Show current migration status")
	fmt.Println("  force <version>- Force set the migration version (use with caution)")
	fmt.Println("  drop           - Drop all tables (DANGEROUS - requires confirmation)")
	fmt.Println("  help           - Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go create \"create_users_table\"")
	fmt.Println("  go run cmd/migrate/main.go migrate")
	fmt.Println("  go run cmd/migrate/main.go rollback")
	fmt.Println("  go run cmd/migrate/main.go fresh")
	fmt.Println("  go run cmd/migrate/main.go purge")
	fmt.Println("  go run cmd/migrate/main.go status")
	fmt.Println("  go run cmd/migrate/main.go force 3")
	fmt.Println("  go run cmd/migrate/main.go drop")
	fmt.Println()
	fmt.Println("Note: Migration files are now in separate .up.sql and .down.sql files")
	fmt.Println("      This follows the golang-migrate library standard")
	fmt.Println()
	fmt.Println("Migration Types:")
	fmt.Println("  - migrate/up:   Apply pending migrations sequentially")
	fmt.Println("  - rollback:     Rollback only the last migration")
	fmt.Println("  - fresh:        Drop all tables and re-apply all migrations (DESTRUCTIVE)")
	fmt.Println("  - purge:        Rollback all migrations to version 0 (preserves tables)")
	fmt.Println("  - drop:         Drop all database tables (DESTRUCTIVE)")
}

// isValidCommand checks if the command is in the list of valid commands
func isValidCommand(command string) bool {
	validCommands := []string{
		"migrate", "up", "rollback", "down", "status", "create",
		"force", "drop", "fresh", "purge", "help", "--help", "-h",
	}
	for _, valid := range validCommands {
		if command == valid {
			return true
		}
	}
	return false
}

// validateMigrationName ensures migration name follows proper naming conventions
func validateMigrationName(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	// Check length
	if len(name) > 50 {
		return fmt.Errorf("migration name cannot exceed 50 characters")
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("migration name can only contain letters, numbers, underscores, and hyphens")
	}

	// Check it doesn't start with number
	if regexp.MustCompile(`^[0-9]`).MatchString(name) {
		return fmt.Errorf("migration name cannot start with a number")
	}

	return nil
}

// validateAndParseVersion validates and parses version number
func validateAndParseVersion(versionStr string) (int, error) {
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, fmt.Errorf("version must be a valid integer")
	}

	if version < 0 {
		return 0, fmt.Errorf("version cannot be negative")
	}

	return version, nil
}

// confirmDangerousOperation asks for user confirmation for destructive operations
func confirmDangerousOperation(operation string) error {
	cfg := config.Get()

	// Skip confirmation in non-production environments if explicitly configured
	if cfg.Environment != "production" {
		fmt.Printf("⚠️  This will %s. Continue? (y/N): ", operation)
	} else {
		fmt.Printf("🚨 PRODUCTION ENVIRONMENT DETECTED! 🚨\n")
		fmt.Printf("This will %s in PRODUCTION!\n", operation)
		fmt.Printf("Type 'YES I UNDERSTAND' to continue: ")
	}

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(response)

	if cfg.Environment == "production" {
		if response != "YES I UNDERSTAND" {
			return fmt.Errorf("operation cancelled - exact confirmation required")
		}
	} else {
		if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
			return fmt.Errorf("operation cancelled by user")
		}
	}

	return nil
}
