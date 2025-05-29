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
		var existingRole model.Role
		result := db.Where("code = ?", role.Code).First(&existingRole)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Role doesn't exist, create it
				if err := db.Create(&role).Error; err != nil {
					return err
				}
				log.Printf("✓ Created role: %s", role.Code)
			} else {
				return result.Error
			}
		} else {
			log.Printf("✓ Role already exists: %s", role.Code)
		}
	}

	return nil
}
