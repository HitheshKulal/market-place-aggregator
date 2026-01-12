package repository

import (
	"context"

	"go-backend-api/internal/models"

	"gorm.io/gorm"
)

type TemplateRepository struct {
	*Repository[models.Template]
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		Repository: NewRepository[models.Template](db),
	}
}

func (r *TemplateRepository) Create(ctx context.Context, template *models.Template) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *TemplateRepository) List(ctx context.Context) ([]*models.Template, error) {
	var templates []*models.Template
	err := r.db.
		WithContext(ctx).
		Find(&templates).
		Error

	return templates, err
}
