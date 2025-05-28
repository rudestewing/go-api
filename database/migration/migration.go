package migration

import (
	"fmt"
	"go-api/config"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	Batch     int       `gorm:"not null"`
	ExecutedAt time.Time `gorm:"autoCreateTime"`
}

type MigrationManager struct {
	db *gorm.DB
}

type MigrationFile struct {
	Name     string
	FilePath string
	UpSQL    string
	DownSQL  string
}

func NewMigrationManager() (*MigrationManager, error) {
	cfg := config.Get()
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	manager := &MigrationManager{db: db}

	// Create migrations table if it doesn't exist
	if err := manager.createMigrationsTable(); err != nil {
		return nil, fmt.Errorf("failed to create migrations table: %w", err)
	}

	return manager, nil
}

func (m *MigrationManager) createMigrationsTable() error {
	// Check if migrations table exists
	var exists bool
	err := m.db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'migrations')").Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		// Create new table
		return m.db.AutoMigrate(&Migration{})
	}

	// Table exists, check if we need to update schema
	var hasIDColumn, hasNameColumn, hasBatchColumn, hasExecutedAtColumn bool

	// Check columns existence
	m.db.Raw("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'migrations' AND column_name = 'id')").Scan(&hasIDColumn)
	m.db.Raw("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'migrations' AND column_name = 'name')").Scan(&hasNameColumn)
	m.db.Raw("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'migrations' AND column_name = 'batch')").Scan(&hasBatchColumn)
	m.db.Raw("SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'migrations' AND column_name = 'executed_at')").Scan(&hasExecutedAtColumn)

	// Add missing columns if needed
	if !hasBatchColumn {
		// Add batch column with default value
		err := m.db.Exec("ALTER TABLE migrations ADD COLUMN batch INTEGER DEFAULT 1").Error
		if err != nil {
			return fmt.Errorf("failed to add batch column: %w", err)
		}

		// Update existing records to have batch = 1
		err = m.db.Exec("UPDATE migrations SET batch = 1 WHERE batch IS NULL").Error
		if err != nil {
			return fmt.Errorf("failed to update batch values: %w", err)
		}

		// Make batch column NOT NULL
		err = m.db.Exec("ALTER TABLE migrations ALTER COLUMN batch SET NOT NULL").Error
		if err != nil {
			return fmt.Errorf("failed to make batch column NOT NULL: %w", err)
		}
	}

	if !hasExecutedAtColumn {
		err := m.db.Exec("ALTER TABLE migrations ADD COLUMN executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP").Error
		if err != nil {
			return fmt.Errorf("failed to add executed_at column: %w", err)
		}
	}

	return nil
}

func (m *MigrationManager) GetExecutedMigrations() ([]string, error) {
	var migrations []Migration
	if err := m.db.Order("executed_at").Find(&migrations).Error; err != nil {
		return nil, err
	}

	executed := make([]string, len(migrations))
	for i, migration := range migrations {
		executed[i] = migration.Name
	}
	return executed, nil
}

func (m *MigrationManager) GetPendingMigrations() ([]MigrationFile, error) {
	executed, err := m.GetExecutedMigrations()
	if err != nil {
		return nil, err
	}

	executedMap := make(map[string]bool)
	for _, name := range executed {
		executedMap[name] = true
	}

	allMigrations, err := m.loadMigrationFiles()
	if err != nil {
		return nil, err
	}

	var pending []MigrationFile
	for _, migration := range allMigrations {
		if !executedMap[migration.Name] {
			pending = append(pending, migration)
		}
	}

	return pending, nil
}

func (m *MigrationManager) loadMigrationFiles() ([]MigrationFile, error) {
	migrationsDir := "database/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return []MigrationFile{}, nil
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []MigrationFile
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		migrationFile, err := m.parseMigrationFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			log.Printf("Warning: Failed to parse migration file %s: %v", file.Name(), err)
			continue
		}

		migrations = append(migrations, migrationFile)
	}

	// Sort migrations by filename (which should include timestamp)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
}

func (m *MigrationManager) parseMigrationFile(filePath string) (MigrationFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return MigrationFile{}, err
	}

	contentStr := string(content)
	parts := strings.Split(contentStr, "-- +migrate Down")

	if len(parts) != 2 {
		return MigrationFile{}, fmt.Errorf("migration file must contain both up and down sections separated by '-- +migrate Down'")
	}

	upSQL := strings.TrimSpace(strings.TrimPrefix(parts[0], "-- +migrate Up"))
	downSQL := strings.TrimSpace(parts[1])

	fileName := filepath.Base(filePath)
	migrationName := strings.TrimSuffix(fileName, ".sql")

	return MigrationFile{
		Name:     migrationName,
		FilePath: filePath,
		UpSQL:    upSQL,
		DownSQL:  downSQL,
	}, nil
}

func (m *MigrationManager) RunMigrations() error {
	pending, err := m.GetPendingMigrations()
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		log.Println("No pending migrations to run")
		return nil
	}

	// Get next batch number
	var lastBatch int
	m.db.Model(&Migration{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch)
	nextBatch := lastBatch + 1

	log.Printf("Running %d pending migrations...", len(pending))

	for _, migration := range pending {
		log.Printf("Running migration: %s", migration.Name)

		if err := m.executeMigration(migration, nextBatch); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.Name, err)
		}

		log.Printf("Migration %s completed successfully", migration.Name)
	}

	log.Println("All migrations completed successfully")
	return nil
}

func (m *MigrationManager) executeMigration(migration MigrationFile, batch int) error {
	// Execute migration in transaction
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Execute the SQL
		sqlDB, err := tx.DB()
		if err != nil {
			return err
		}

		if _, err := sqlDB.Exec(migration.UpSQL); err != nil {
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}

		// Record migration
		migrationRecord := Migration{
			Name:  migration.Name,
			Batch: batch,
		}

		if err := tx.Create(&migrationRecord).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})
}

func (m *MigrationManager) RollbackLastBatch() error {
	// Get last batch
	var lastBatch int
	if err := m.db.Model(&Migration{}).Select("MAX(batch)").Scan(&lastBatch).Error; err != nil {
		return err
	}

	if lastBatch == 0 {
		log.Println("No migrations to rollback")
		return nil
	}

	// Get migrations from last batch
	var migrations []Migration
	if err := m.db.Where("batch = ?", lastBatch).Order("executed_at DESC").Find(&migrations).Error; err != nil {
		return err
	}

	log.Printf("Rolling back %d migrations from batch %d...", len(migrations), lastBatch)

	for _, migration := range migrations {
		log.Printf("Rolling back migration: %s", migration.Name)

		if err := m.rollbackMigration(migration); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.Name, err)
		}

		log.Printf("Migration %s rolled back successfully", migration.Name)
	}

	log.Println("Rollback completed successfully")
	return nil
}

func (m *MigrationManager) rollbackMigration(migration Migration) error {
	// Load migration file to get down SQL
	migrationFiles, err := m.loadMigrationFiles()
	if err != nil {
		return err
	}

	var migrationFile *MigrationFile
	for _, file := range migrationFiles {
		if file.Name == migration.Name {
			migrationFile = &file
			break
		}
	}

	if migrationFile == nil {
		return fmt.Errorf("migration file not found for: %s", migration.Name)
	}

	// Execute rollback in transaction
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Execute the down SQL
		sqlDB, err := tx.DB()
		if err != nil {
			return err
		}

		if _, err := sqlDB.Exec(migrationFile.DownSQL); err != nil {
			return fmt.Errorf("failed to execute rollback SQL: %w", err)
		}

		// Remove migration record
		if err := tx.Delete(&migration).Error; err != nil {
			return fmt.Errorf("failed to remove migration record: %w", err)
		}

		return nil
	})
}

func (m *MigrationManager) GetMigrationStatus() error {
	executed, err := m.GetExecutedMigrations()
	if err != nil {
		return err
	}

	pending, err := m.GetPendingMigrations()
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	if len(executed) > 0 {
		fmt.Println("\nExecuted migrations:")
		for _, name := range executed {
			fmt.Printf("  ✓ %s\n", name)
		}
	}

	if len(pending) > 0 {
		fmt.Println("\nPending migrations:")
		for _, migration := range pending {
			fmt.Printf("  ⏳ %s\n", migration.Name)
		}
	} else {
		fmt.Println("\nNo pending migrations")
	}

	return nil
}

func (m *MigrationManager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
