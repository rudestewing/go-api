package cmd

import (
	"fmt"
	"go-api/database/seeder"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Database seeder management",
	Long: `Database seeder management tool for creating and running seeders.

Examples:
  go-api seed create users_seeder
  go-api seed run database/seeders/20250529000000_roles.go`,
}

// seedCreateCmd represents the seed create command
var seedCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new seeder file",
	Long: `Create a new seeder file with the given name.

The seeder name should follow these rules:
- Only contain letters, numbers, underscores, and hyphens
- Cannot start with a number
- Cannot exceed 50 characters

Example:
  go-api seed create users_seeder`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		
		if err := validateSeederName(name); err != nil {
			log.Fatalf("Invalid seeder name: %v", err)
		}
		
		if err := seeder.CreateSeeder(name); err != nil {
			log.Fatalf("Failed to create seeder: %v", err)
		}
	},
}

// seedRunCmd represents the seed run command
var seedRunCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Run a specific seeder file",
	Long: `Run a specific seeder file by providing its path.

The seeder file must be:
- Located in the database/seeders/ directory
- Have a .go extension
- Be a valid Go file with proper seeder implementation

Example:
  go-api seed run database/seeders/20250529000000_roles.go`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		
		if err := validateSeederPath(filePath); err != nil {
			log.Fatalf("Invalid seeder path: %v", err)
		}
		
		if err := seeder.RunSeederFile(filePath); err != nil {
			log.Fatalf("Failed to run seeder file: %v", err)
		}
	},
}

// validateSeederName ensures seeder name follows proper naming conventions
func validateSeederName(name string) error {
	if name == "" {
		return fmt.Errorf("seeder name cannot be empty")
	}

	if len(name) > 50 {
		return fmt.Errorf("seeder name cannot exceed 50 characters")
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("seeder name can only contain letters, numbers, underscores, and hyphens")
	}

	if regexp.MustCompile(`^[0-9]`).MatchString(name) {
		return fmt.Errorf("seeder name cannot start with a number")
	}

	return nil
}

// validateSeederPath ensures seeder file path is secure and valid
func validateSeederPath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("seeder file path cannot be empty")
	}

	// Ensure path is within allowed directories
	cleanPath := filepath.Clean(filePath)
	// Use filepath.ToSlash to normalize path separators for cross-platform compatibility
	normalizedPath := filepath.ToSlash(cleanPath)
	if !strings.HasPrefix(normalizedPath, "database/seeders/") {
		return fmt.Errorf("seeder files must be in database/seeders/ directory")
	}

	// Ensure it's a Go file
	if !strings.HasSuffix(cleanPath, ".go") {
		return fmt.Errorf("seeder files must have .go extension")
	}

	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	return nil
}

func init() {
	// Add seed command to root command
	RootCmd.AddCommand(seedCmd)
	
	// Add subcommands to seed command
	seedCmd.AddCommand(seedCreateCmd)
	seedCmd.AddCommand(seedRunCmd)
}