package main

import (
	"fmt"
	"go-api/database/seeder"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "help", "--help", "-h":
		showUsage()
		return
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Seeder name is required. Usage: go run cmd/seed/main.go create \"seeder_name\"")
		}
		createSeeder(os.Args[2])
		return
	case "run":
		if len(os.Args) < 3 {
			log.Fatal("Seeder file path is required. Usage: go run cmd/seed/main.go run \"path/to/seeder.go\"")
		}
		runSeederFile(os.Args[2])
		return
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showUsage()
	}
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
