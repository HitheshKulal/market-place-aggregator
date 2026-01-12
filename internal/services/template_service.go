package services

import (
	"go-backend-api/internal/api/requests"
	"go-backend-api/internal/models"
	"go-backend-api/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
)

type TemplateServiceConfig struct {
	Repositories *repository.Repositories
	Services     *Services
}

type TemplateService struct {
	repo     *repository.TemplateRepository
	services *Services
}

func NewTemplateService(config *TemplateServiceConfig) *TemplateService {
	return &TemplateService{
		repo:     config.Repositories.Template,
		services: config.Services,
	}
}

func (s *TemplateService) Store(c *gin.Context, req requests.CreateTemplateRequest) (*models.Template, error) {
	template := &models.Template{
		Name:      req.Name,
		Fields:    req.Fields,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(c.Request.Context(), template)
	if err != nil {
		return nil, err
	}

	return template, err
}

func (s *TemplateService) Index(c *gin.Context) ([]*models.Template, error) {
	return s.repo.List(c.Request.Context())
}
