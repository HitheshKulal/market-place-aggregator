package services

import (
	"go-backend-api/internal/repository"
)

type Services struct {
	Product  *ProductService
	Template *TemplateService
	Mapping  *MappingService
	Seller   *SellerService
}

type ServiceConfig struct {
	Repositories *repository.Repositories
}

func NewServices(config *ServiceConfig) *Services {
	services := &Services{}

	services.Product = NewProductService(&ProductServiceConfig{
		Repositories: config.Repositories,
		Services:     services,
	})

	services.Template = NewTemplateService(&TemplateServiceConfig{
		Repositories: config.Repositories,
		Services:     services,
	})

	services.Mapping = NewMappingService(&MappingServiceConfig{
		Repositories: config.Repositories,
		Services:     services,
	})

	services.Seller = NewSellerService(&SellerServiceConfig{
		Repositories: config.Repositories,
		Services:     services,
	})

	return services
}
