package services

import (
	"context"
	"go-backend-api/internal/models"
	"go-backend-api/internal/repository"
)

type SellerServiceConfig struct {
	Repositories *repository.Repositories
	Services     *Services
}

type SellerService struct {
	repo     *repository.SellerRepository
	services *Services
}

func NewSellerService(config *SellerServiceConfig) *SellerService {
	return &SellerService{
		repo:     config.Repositories.Seller,
		services: config.Services,
	}
}

func (s *SellerService) Create(ctx context.Context, sellerModel *models.Seller) error {
	return s.repo.Create(ctx, sellerModel)
}
