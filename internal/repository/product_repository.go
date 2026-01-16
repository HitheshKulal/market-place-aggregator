package repository

import (
	"context"

	"go-backend-api/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	*Repository[models.Product]
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		Repository: NewRepository[models.Product](db),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) GetBySellerID(ctx context.Context, sellerID uint) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Find(&products).
		Where("seller_id = ?", sellerID).
		Error
	return products, err
}

func (r *ProductRepository) List(ctx context.Context) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Find(&products).
		Error
	return products, err
}
