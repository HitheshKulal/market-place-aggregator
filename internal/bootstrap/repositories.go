package bootstrap

import (
	"go-backend-api/internal/logger"
	"go-backend-api/internal/repository"

	"gorm.io/gorm"
)

func InitRepositories(db *gorm.DB) *repository.Repositories {
	repos := repository.NewRepositories(db)
	logger.Info("Repositories initialized successfully")
	return repos
}
