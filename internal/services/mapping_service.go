package services

import (
	"context"
	"go-backend-api/internal/api/requests"
	"go-backend-api/internal/models"
	"go-backend-api/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
)

type MappingServiceConfig struct {
	Repositories *repository.Repositories
	Services     *Services
}

type MappingService struct {
	repo     *repository.MappingRepository
	services *Services
}

func NewMappingService(config *MappingServiceConfig) *MappingService {
	return &MappingService{
		repo:     config.Repositories.Mapping,
		services: config.Services,
	}
}

func (s *MappingService) FindByID(ctx context.Context, templateID uint) (*models.Mapping, error) {
	return s.repo.FindByID(ctx, templateID)
}

func (s *MappingService) Index(c *gin.Context) ([]*models.Mapping, error) {
	return s.repo.List(c.Request.Context())
}

func (s *MappingService) Store(c *gin.Context, req requests.CreateMappingRequest) (*models.Mapping, error) {
	template := &models.Mapping{
		TemplateID: req.TemplateID,
		SellerID:   req.SellerID,
		FieldMap:   req.FieldMap,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := s.repo.Create(c.Request.Context(), template)
	if err != nil {
		return nil, err
	}

	return template, err
}
