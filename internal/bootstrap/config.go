package bootstrap

import (
	"go-backend-api/internal/config"
	"go-backend-api/internal/logger"

	"go.uber.org/zap"
)

func InitConfig() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}
	return cfg
}
