package bootstrap

import (
	"go-backend-api/internal/logger"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	if err := godotenv.Load(); err != nil {
		if env == "development" {
			logger.Fatal("Error loading .env file", zap.Error(err))
		}
		logger.Info("No .env file found, using environment variables")
	}
}
