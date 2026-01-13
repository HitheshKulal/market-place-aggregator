package repository

import (
	"context"

	"go-backend-api/internal/models"

	"gorm.io/gorm"
)

type MappingRepository struct {
	*Repository[models.Mapping]
}

func NewMappingRepository(db *gorm.DB) *MappingRepository {
	return &MappingRepository{
		Repository: NewRepository[models.Mapping](db),
	}
}

func (r *MappingRepository) Create(ctx context.Context, mapping *models.Mapping) error {
	return r.db.WithContext(ctx).Create(mapping).Error
}

func (r *MappingRepository) FindByID(ctx context.Context, id uint) (*models.Mapping, error) {
	var mapping *models.Mapping
	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&mapping).
		Error
	return mapping, err
}

func (r *MappingRepository) List(ctx context.Context) ([]*models.Mapping, error) {
	var mappings []*models.Mapping
	err := r.db.
		WithContext(ctx).
		Find(&mappings).
		Error
	return mappings, err
}
