package main

import (
	"fmt"
	"go-api/database/seeder"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	command := os.Args[1]

	// Validate command
	if !isValidSeederCommand(command) {
		fmt.Printf("Unknown command: %s\n\n", command)
		showUsage()
		return
	}

	switch command {
	case "help", "--help", "-h":
		showUsage()
		return
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Seeder name is required. Usage: go run cmd/seed/main.go create \"seeder_name\"")
		}
		if err := validateSeederName(os.Args[2]); err != nil {
			log.Fatalf("Invalid seeder name: %v", err)
		}
		createSeeder(os.Args[2])
		return
	case "run":
		if len(os.Args) < 3 {
			log.Fatal("Seeder file path is required. Usage: go run cmd/seed/main.go run \"path/to/seeder.go\"")
		}
		if err := validateSeederPath(os.Args[2]); err != nil {
			log.Fatalf("Invalid seeder path: %v", err)
		}
		runSeederFile(os.Args[2])
		return
	}
}

// isValidSeederCommand checks if the command is valid
func isValidSeederCommand(command string) bool {
	validCommands := []string{"help", "--help", "-h", "create", "run"}
	for _, valid := range validCommands {
		if command == valid {
			return true
		}
	}
	return false
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
	if !strings.HasPrefix(cleanPath, "database/seeders/") {
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

func runSeederFile(filePath string) {
	if err := seeder.RunSeederFile(filePath); err != nil {
		log.Fatalf("Failed to run seeder file: %v", err)
	}
}

func createSeeder(name string) {
	if err := seeder.CreateSeeder(name); err != nil {
		log.Fatalf("Failed to create seeder: %v", err)
	}
}

func showUsage() {
	fmt.Println("Database Seeder Tool")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed/main.go <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create <name>      - Create a new seeder file (.go)")
	fmt.Println("  run <path>         - Run a specific seeder file")
	fmt.Println("  help               - Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed/main.go create \"users_seeder\"")
	fmt.Println("  go run cmd/seed/main.go run \"database/seeders/20250529000000_roles.go\"")
	fmt.Println()
	fmt.Println("Note: Seeder files are .go files in the database/seeders directory")
	fmt.Println("      Each seeder should implement Run() method")
}
