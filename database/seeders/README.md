# Database Seeder Feature

This project includes a simplified database seeder feature. Seeders allow you to populate your database with test data using individual Go files.

## Features

- Create new seeder files with timestamps
- Run individual seeder files manually
- Each seeder is a standalone Go program
- No batch tracking or rollback functionality (simplified approach)

## Usage

### Create a New Seeder

```bash
make seed-create name="users_seeder"
# or
go run cmd/seed/main.go create "users_seeder"
```

This will create a new seeder file in `database/seeders/` with a timestamp prefix.

### Run a Specific Seeder

```bash
make seed-run path="database/seeders/20250529000000_roles.go"
# or
go run cmd/seed/main.go run "database/seeders/20250529000000_roles.go"
```

This will execute the specific seeder file.

### Show Help

```bash
make seed-help
# or
go run cmd/seed/main.go help
```

- List of all executed seeders
- List of available registered seeders

### Show Help

```bash
go run cmd/seed/main.go help
```

## How Seeders Work

### Seeder Structure

Each seeder is a standalone Go program with a `main` function that:

1. Initializes the database configuration
2. Connects to the database
3. Runs the seeding logic
4. Handles errors appropriately

### Example Seeder

When you create a seeder using the `create` command, it generates a template like this:

```go
package main

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

    // TODO: Add your seeding logic here
    // Example:
    // users := []model.User{
    //     {Name: "John Doe", Email: "john@example.com"},
    //     {Name: "Jane Smith", Email: "jane@example.com"},
    // }
    //
    // for _, user := range users {
    //     var existingUser model.User
    //     result := db.Where("email = ?", user.Email).First(&existingUser)
    //     if result.Error != nil {
    //         if result.Error == gorm.ErrRecordNotFound {
    //             if err := db.Create(&user).Error; err != nil {
    //                 return err
    //             }
    //             log.Printf("✓ Created user: %s", user.Email)
    //         } else {
    //             return result.Error
    //         }
    //     } else {
    //         log.Printf("✓ User already exists: %s", user.Email)
    //     }
    // }

    return nil
}
```

// Rollback reverses the seeder
func (s *UsersSeeder) Rollback(db *gorm.DB) error {
log.Printf("Rolling back users seeder...")

    // TODO: Add your rollback logic here
    // Example:
    // if err := db.Where("email IN ?", []string{"john@example.com", "jane@example.com"}).Delete(&model.User{}).Error; err != nil {
    //     return err
    // }

    log.Printf("✓ users seeder rollback completed")
    return nil

}

// GetName returns the seeder name
func (s \*UsersSeeder) GetName() string {
return "users"
}

// init function to auto-register the seeder
func init() {
seeder.RegisterSeeder("20250529000001_users", &UsersSeeder{})
}

```

### 3. Auto-Registration

## Best Practices

1. **Use descriptive names**: Name your seeders clearly to indicate what data they seed.
2. **Handle duplicates**: Always check if data already exists before creating new records.
3. **Use transactions**: Wrap your seeding logic in database transactions for data integrity.
4. **Keep seeders idempotent**: Seeders should be safe to run multiple times.

## Directory Structure

```

database/
├── seeder/
│ └── seeder.go # Seeder utilities and functions
├── seeders/
│ ├── README.md # This documentation
│ ├── roles/
│ │ └── 20250529000000_roles.go # Roles seeder
│ └── users/
│ └── 20250529000001_users.go # Users seeder
└── ...

cmd/
├── seed/
│ └── main.go # Seeder command-line tool
└── ...

```

## Notes

- Each seeder is a standalone Go program that can be run independently
- Seeders are organized in subdirectories or as individual files
- No automatic registration or batch tracking (simplified approach)
- Use `make seed-run path="path/to/seeder.go"` to run specific seeders

## Example Seeders

Two example seeders are included:

1. **Roles Seeder** (`roles/20250529000000_roles.go`): Seeds basic roles (admin, user, moderator)
2. **Users Seeder** (`users/20250529000001_users.go`): Seeds example users with roles
```
