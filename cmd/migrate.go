package cmd

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

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration management",
	Long: `Database migration management tool for creating and running migrations.

This tool uses golang-migrate library standard for SQL migrations with
separate .up.sql and .down.sql files.

Examples:
  go-api migrate up
  go-api migrate create create_users_table
  go-api migrate rollback`,
}

// migrateUpCmd represents the migrate up command
var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Long:  `Run all pending migrations to bring the database schema up to date.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		runMigrations()
	},
}

// migrateDownCmd represents the migrate down command
var migrateDownCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"rollback"},
	Short:   "Rollback the last migration",
	Long:    `Rollback the last applied migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		rollbackMigrations()
	},
}

// migrateStatusCmd represents the migrate status command
var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current migration status",
	Long:  `Display the current migration status including applied and pending migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		showStatus()
	},
}

// migrateCreateCmd represents the migrate create command
var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file pair",
	Long: `Create a new migration file pair (.up.sql and .down.sql) with the given name.

The migration name should follow these rules:
- Only contain letters, numbers, underscores, and hyphens
- Cannot start with a number
- Cannot exceed 50 characters

Example:
  go-api migrate create create_users_table`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		
		if err := validateMigrationName(name); err != nil {
			log.Fatalf("Invalid migration name: %v", err)
		}
		
		config.InitConfig()
		createMigration(name)
	},
}

// migrateForceCmd represents the migrate force command
var migrateForceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force set the migration version (use with caution)",
	Long: `Force set the migration version to a specific number.
This is a dangerous operation that should only be used to fix migration state issues.

Example:
  go-api migrate force 3`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		versionStr := args[0]
		
		version, err := validateAndParseVersion(versionStr)
		if err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		
		config.InitConfig()
		forceVersion(version)
	},
}

// migrateDropCmd represents the migrate drop command
var migrateDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop all tables (DANGEROUS)",
	Long: `Drop all tables in the database.
This is a destructive operation that requires confirmation.

WARNING: This will permanently delete all data in your database!`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		
		if err := confirmDangerousOperation("drop all database tables"); err != nil {
			log.Fatalf("Operation cancelled: %v", err)
		}
		
		dropDatabase()
	},
}

// migrateFreshCmd represents the migrate fresh command
var migrateFreshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Long: `Drop all tables and re-run all migrations from the beginning.
This is a destructive operation that requires confirmation.

WARNING: This will permanently delete all data in your database!`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		
		if err := confirmDangerousOperation("drop and recreate all database tables"); err != nil {
			log.Fatalf("Operation cancelled: %v", err)
		}
		
		runFreshMigrations()
	},
}

// migratePurgeCmd represents the migrate purge command
var migratePurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Rollback all executed migrations to version 0",
	Long: `Rollback all executed migrations to version 0.
This preserves table structure but marks all migrations as not applied.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		
		if err := confirmDangerousOperation("purge all migrations"); err != nil {
			log.Fatalf("Operation cancelled: %v", err)
		}
		
		purgeMigrations()
	},
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

	if err := manager.RollbackLastMigration(); err != nil {
		log.Fatalf("Failed to rollback last migration: %v", err)
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
	manager, err := migration.NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer manager.Close()

	if err := manager.Purge(); err != nil {
		log.Fatalf("Failed to purge migrations: %v", err)
	}
}

// validateMigrationName ensures migration name follows proper naming conventions
func validateMigrationName(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	if len(name) > 50 {
		return fmt.Errorf("migration name cannot exceed 50 characters")
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("migration name can only contain letters, numbers, underscores, and hyphens")
	}

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

func init() {
	// Add migrate command to root command
	RootCmd.AddCommand(migrateCmd)
	
	// Add subcommands to migrate command
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateForceCmd)
	migrateCmd.AddCommand(migrateDropCmd)
	migrateCmd.AddCommand(migrateFreshCmd)
	migrateCmd.AddCommand(migratePurgeCmd)
}