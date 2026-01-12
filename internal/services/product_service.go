package services

import (
	"context"
	"errors"
	"go-backend-api/internal/repository"
)

type ProductServiceConfig struct {
	Repositories *repository.Repositories
	Services     *Services
}

type ProductService struct {
	repo     *repository.ProductRepository
	services *Services
}

func NewProductService(config *ProductServiceConfig) *ProductService {
	return &ProductService{
		repo:     config.Repositories.Product,
		services: config.Services,
	}
}

// GetProductsByTemplate transforms all products according to a specific template mapping
// Input: templateID
// Output: list of products with data transformed according to template format
// Example: [{"name":"Iphone", "price":551511}] -> [{"title":"Iphone", "price":551511}]
func (s *ProductService) GetProductsByTemplate(ctx context.Context, templateID uint) ([]map[string]interface{}, error) {
	// Get mapping for this template
	mapping, err := s.services.Mapping.FindByTemplateID(ctx, templateID)
	if err != nil {
		return nil, errors.New("no mapping found for template")
	}

	// Get all products
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Transform all products according to mapping
	result := make([]map[string]interface{}, 0, len(products))
	for _, product := range products {
		transformedProduct := s.transformProduct(product.Data, mapping.FieldMap)
		result = append(result, transformedProduct)
	}

	return result, nil
}

// transformProduct applies field mapping to a single product
// fieldMap can be either {"title":"name"} or {"name":"title"} depending on direction
func (s *ProductService) transformProduct(productData map[string]interface{}, fieldMap map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	// fieldMap format: {"templateField":"productField"}
	// e.g., {"title":"name", "price":"price"}
	for templateField, productField := range fieldMap {
		if value, exists := productData[productField]; exists {
			result[templateField] = value
		}
	}

	return result
}
