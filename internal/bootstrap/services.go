package bootstrap

import (
	"go-backend-api/internal/logger"
	"go-backend-api/internal/repository"
	"go-backend-api/internal/services"
)

func InitServices(
	repos *repository.Repositories,
) *services.Services {
	servicesCfg := &services.ServiceConfig{
		Repositories: repos,
	}
	svc := services.NewServices(servicesCfg)
	logger.Info("Services initialized successfully")
	return svc
}
