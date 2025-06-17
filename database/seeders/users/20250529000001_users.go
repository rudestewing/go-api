package main

import (
	"crypto/rand"
	"go-api/config"
	"go-api/database"
	"go-api/model"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
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

	// Generate secure passwords for seeded users
	adminPassword, err := generateSecurePassword()
	if err != nil {
		return err
	}

	userPassword, err := generateSecurePassword()
	if err != nil {
		return err
	}

	adminHashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash admin password:", err)
	}

	userHashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash user password:", err)
	}

	// Example users data
	users := []model.User{
		{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: string(adminHashedPassword),
			RoleID:   1, // Assuming admin role
		},
		{
			Name:     "Regular User",
			Email:    "user@example.com",
			Password: string(userHashedPassword),
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

	// Log generated passwords for development (remove in production)
	log.Printf("Generated passwords (for development only):")
	log.Printf("Admin password: %s", adminPassword)
	log.Printf("User password: %s", userPassword)

	return nil
}

// generateSecurePassword generates a cryptographically secure random password
func generateSecurePassword() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
	const passwordLength = 16

	password := make([]byte, passwordLength)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}
	return string(password), nil
}
