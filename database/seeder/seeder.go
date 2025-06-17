package seeder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// SeederInterface defines the interface for seeders
type SeederInterface interface {
	Run(db *gorm.DB) error
	GetName() string
}

// RunSeederFile runs a specific seeder file by compiling and executing it
func RunSeederFile(filePath string) error {
	// Security validation - ensure seeder files are in the correct directory
	if !strings.HasPrefix(filePath, "database/seeders/") {
		return fmt.Errorf("seeder files must be in database/seeders/ directory")
	}

	if !strings.HasSuffix(filePath, ".go") {
		return fmt.Errorf("seeder files must have .go extension")
	}

	// Check if file exists and validate it's readable
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("seeder file does not exist: %s", filePath)
	}
	if err != nil {
		return fmt.Errorf("cannot access seeder file %s: %w", filePath, err)
	}

	// Ensure it's a regular file, not a directory or symlink
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("seeder path must be a regular file: %s", filePath)
	}

	log.Printf("Running seeder file: %s", filePath)

	// For Go seeder files, we can use go run to execute them
	// This is a simple approach that compiles and runs the seeder
	cmd := exec.Command("go", "run", filePath)
	cmd.Dir = "." // Current working directory

	// Set environment variables if needed
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute seeder file %s: %w\nOutput: %s", filePath, err, string(output))
	}

	log.Printf("✓ Seeder file %s executed successfully", filepath.Base(filePath))
	if len(output) > 0 {
		log.Printf("Output: %s", string(output))
	}

	return nil
}

// CreateSeeder creates a new seeder file
func CreateSeeder(name string) error {
	if name == "" {
		return fmt.Errorf("seeder name is required")
	}

	// Create seeders directory if it doesn't exist
	seedersDir := "database/seeders"
	if err := os.MkdirAll(seedersDir, 0750); err != nil {
		return fmt.Errorf("failed to create seeders directory: %w", err)
	}

	// Generate timestamp
	timestamp := time.Now().Format("20060102150405")

	// Clean seeder name (replace spaces with underscores, remove special chars)
	cleanName := strings.ReplaceAll(name, " ", "_")
	cleanName = strings.ToLower(cleanName)

	// Title case for struct name
	titleName := strings.ToUpper(cleanName[:1]) + cleanName[1:]

	// Create filename
	filename := fmt.Sprintf("%s_%s.go", timestamp, cleanName)
	filepath := filepath.Join(seedersDir, filename)

	// Create seeder file template
	template := fmt.Sprintf(`package main

import (
	"go-api/config"
	"go-api/database"
	"log"

	"gorm.io/gorm"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %%v", err)
	}

	// Run the seeder
	if err := run%s(db); err != nil {
		log.Fatalf("Failed to run %s seeder: %%v", err)
	}

	log.Printf("✓ %s seeder completed successfully")
}

func run%s(db *gorm.DB) error {
	log.Printf("Running %s seeder...")

	// TODO: Add your seeding logic here
	// Example:
	// users := []model.User{
	//     {Name: "John Doe", Email: "john@example.com"},
	//     {Name: "Jane Smith", Email: "jane@example.com"},
	// }
	//
	// for _, user := range users {
	//     if err := db.FirstOrCreate(&user, model.User{Email: user.Email}).Error; err != nil {
	//         return err
	//     }
	// }

	return nil
}
`,
		titleName, cleanName, cleanName,
		titleName, cleanName)

	// Write seeder file
	if err := os.WriteFile(filepath, []byte(template), 0600); err != nil {
		return fmt.Errorf("failed to create seeder file: %w", err)
	}

	fmt.Printf("Seeder file created:\n")
	fmt.Printf("  📝 %s\n", filepath)
	fmt.Println("Please edit the file to add your seeding logic")
	fmt.Printf("Run it with: make seed-run path=\"%s\"\n", filepath)

	return nil
}
