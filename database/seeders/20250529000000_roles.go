package main

import (
	"go-api/app/model"
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
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run the seeder
	if err := runRoles(db); err != nil {
		log.Fatalf("Failed to run roles seeder: %v", err)
	}

	log.Printf("✓ roles seeder completed successfully")
}

func runRoles(db *gorm.DB) error {
	log.Printf("Running roles seeder...")

	// Example roles data
	roles := []model.Role{
		{
			Code: "admin",
			Name: "Administrator",
		},
		{
			Code: "user",
			Name: "User",
		},
		{
			Code: "moderator",
			Name: "Moderator",
		},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Code: role.Code}).Error; err != nil {
			return err
		}
	}

	return nil
}
