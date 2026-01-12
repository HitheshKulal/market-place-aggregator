package bootstrap

import (
	"go-backend-api/internal/config"
	"go-backend-api/internal/db"
	"go-backend-api/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dbInstance, err := db.NewDatabase(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize database",
			zap.Error(err),
			zap.String("host", cfg.Database.Host),
			zap.Int("port", cfg.Database.Port),
			zap.String("database", cfg.Database.Name),
		)
	}
	sqlDB, err := dbInstance.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", zap.Error(err))
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("Database connection established successfully")
	return dbInstance
}
