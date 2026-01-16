package repository

import (
	"context"

	"go-backend-api/internal/models"

	"gorm.io/gorm"
)

type SellerRepository struct {
	*Repository[models.Seller]
}

func NewSellerRepository(db *gorm.DB) *SellerRepository {
	return &SellerRepository{
		Repository: NewRepository[models.Seller](db),
	}
}

func (r *SellerRepository) Create(ctx context.Context, seller *models.Seller) error {
	return r.db.WithContext(ctx).Create(seller).Error
}
