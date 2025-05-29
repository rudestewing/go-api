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
	if err := runUsers(db); err != nil {
		log.Fatalf("Failed to run users seeder: %v", err)
	}

	log.Printf("✓ users seeder completed successfully")
}

func runUsers(db *gorm.DB) error {
	log.Printf("Running users seeder...")

	// Example users data
	users := []model.User{
		{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			RoleID:   1, // Assuming admin role
		},
		{
			Name:     "Regular User",
			Email:    "user@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			RoleID:   2, // Assuming user role
		},
	}

	for _, user := range users {
		var existingUser model.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// User doesn't exist, create it
				if err := db.Create(&user).Error; err != nil {
					return err
				}
				log.Printf("✓ Created user: %s", user.Email)
			} else {
				return result.Error
			}
		} else {
			log.Printf("✓ User already exists: %s", user.Email)
		}
	}

	return nil
}
