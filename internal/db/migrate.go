package db

import (
	"fmt"
	"go-backend-api/internal/logger"
	"go-backend-api/internal/models"

	"time"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	migrations := []struct {
		name     string
		migrator func(*gorm.DB) error
	}{
		{
			name: "create_product_table",
			migrator: func(db *gorm.DB) error {
				return models.AutoMigrateProduct(db)
			},
		},
		{
			name: "crate_mapping_table",
			migrator: func(db *gorm.DB) error {
				return models.AutoMigrateMapping(db)
			},
		},
		{
			name: "create_template_table",
			migrator: func(db *gorm.DB) error {
				return models.AutoMigrateTemplate(db)
			},
		},
	}

	for _, migration := range migrations {
		var existing Migration
		if err := db.Where("name = ?", migration.name).First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return fmt.Errorf("failed to check migration status: %w", err)
			}

			if err := migration.migrator(db); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.name, err)
			}

			if err := db.Create(&Migration{
				Name:      migration.name,
				AppliedAt: time.Now(),
			}).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
			}
		} else {
			fmt.Printf("Migration %s already applied at %s\n", migration.name, existing.AppliedAt)
		}
	}

	logger.GetLogger().Info("------------------------Migration Completed--------------")

	return nil
}
