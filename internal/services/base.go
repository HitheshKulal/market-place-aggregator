package services

import (
	"go-backend-api/internal/repository"
)

type Services struct {
}

type ServiceConfig struct {
	Repositories *repository.Repositories
}

func NewServices(config *ServiceConfig) *Services {
	services := &Services{}

	return services
}
