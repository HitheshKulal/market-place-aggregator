package services

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"go-backend-api/internal/api/requests"
	"go-backend-api/internal/models"
	"go-backend-api/internal/repository"

	"mime/multipart"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
func (s *ProductService) GetProductsByTemplate(ctx context.Context, mappingID uint) ([]map[string]interface{}, error) {
	// Get mapping for this template
	mapping, err := s.services.Mapping.FindByID(ctx, mappingID)
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
		// Convert product struct to map first
		productMap := s.productToMap(product)

		// Then transform according to template mapping
		transformedProduct := s.transformProduct(productMap, mapping.FieldMap)
		result = append(result, transformedProduct)
	}

	return result, nil
}

// productToMap converts Product struct fields to a map
func (s *ProductService) productToMap(product models.Product) map[string]interface{} {
	productMap := make(map[string]interface{})

	productMap["id"] = product.ID
	productMap["sku"] = product.SKU
	productMap["name"] = product.Name
	productMap["brand_name"] = product.BrandName
	productMap["gender"] = product.Gender
	productMap["category"] = product.Category
	productMap["color"] = product.Color
	productMap["size"] = product.Size
	productMap["material"] = product.Material
	productMap["image1"] = product.Image1
	productMap["image2"] = product.Image2
	productMap["description"] = product.Description
	productMap["mrp"] = product.MRP
	productMap["price"] = product.Price
	productMap["quantity"] = product.Quantity

	// Merge additional data from JSONB field
	for key, value := range product.Data {
		// Don't override existing mapped fields
		if _, exists := productMap[key]; !exists {
			productMap[key] = value
		}
	}

	return productMap
}

// transformProduct applies field mapping to a single product
func (s *ProductService) transformProduct(productData map[string]interface{}, fieldMap map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	// fieldMap format: {"templateField":"productField"}
	// e.g., {"productName":"name", "brand":"brand_name"}
	for templateField, productField := range fieldMap {
		if value, exists := productData[productField]; exists {
			result[templateField] = value
		}
	}

	return result
}

func (s *ProductService) UploadAndStoreProducts(c *gin.Context, file multipart.File, filename string) (*requests.UploadProductsResponse, error) {
	// Step 1: Extract CSV data
	headers, dataRows, err := s.extractCSVData(file)
	if err != nil {
		return nil, err
	}

	response := &requests.UploadProductsResponse{
		FileName:          filename,
		TotalRows:         len(dataRows),
		DiscoveredColumns: headers,
		SampleProducts:    make([]map[string]interface{}, 0),
		Errors:            make([]requests.ProductError, 0),
	}

	// Step 2: Map and validate each row
	validatedProducts := make([]*models.Product, 0)
	productMaps := make([]map[string]interface{}, 0)

	for i, row := range dataRows {
		rowNumber := i + 2 // +2 because: +1 for 0-index, +1 for header row

		// Map CSV row to product data
		productData := s.mapRowToProductData(headers, row)

		// Validate and create product
		product, err := s.validateAndCreateProduct(productData)
		if err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, requests.ProductError{
				Row:     rowNumber,
				SKU:     getString(productData, "sku"),
				Message: err.Error(),
			})
			continue
		}

		validatedProducts = append(validatedProducts, product)
		productMaps = append(productMaps, productData)
	}

	// Step 3: Store validated products
	var createdProducts []map[string]interface{}

	for i, product := range validatedProducts {
		err := s.repo.Create(c.Request.Context(), product)
		if err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, requests.ProductError{
				Row:     i + 2,
				SKU:     product.SKU,
				Message: fmt.Sprintf("Database error: %v", err),
			})
			continue
		}

		response.SuccessCount++

		// Add to sample (first 5 successful products)
		if len(createdProducts) < 5 {
			createdProducts = append(createdProducts, map[string]interface{}{
				"id":       product.ID,
				"sku":      product.SKU,
				"name":     product.Name,
				"price":    product.Price,
				"quantity": product.Quantity,
			})
		}
	}

	response.SampleProducts = createdProducts

	// If all products failed, return error
	if response.FailedCount == response.TotalRows {
		return response, errors.New("all products failed to import")
	}

	return response, nil
}

// extractCSVData - Step 1: Extract and parse CSV file
func (s *ProductService) extractCSVData(file multipart.File) ([]string, [][]string, error) {
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, nil, errors.New("CSV file is empty")
	}

	if len(records) < 2 {
		return nil, nil, errors.New("CSV file must have at least headers and one data row")
	}

	// First row is headers
	headers := records[0]
	// Trim whitespace from headers
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	// Validate that headers are not empty
	for _, header := range headers {
		if header == "" {
			return nil, nil, errors.New("CSV contains empty column headers")
		}
	}

	dataRows := records[1:]
	return headers, dataRows, nil
}

// mapRowToProductData - Step 2a: Map CSV row to product data structure
func (s *ProductService) mapRowToProductData(headers []string, row []string) map[string]interface{} {
	productData := make(map[string]interface{})
	for j, value := range row {
		if j < len(headers) {
			productData[headers[j]] = strings.TrimSpace(value)
		}
	}
	return productData
}

// validateAndCreateProduct - Step 2b: Validate product data and create product model
func (s *ProductService) validateAndCreateProduct(productData map[string]interface{}) (*models.Product, error) {
	product := &models.Product{
		SKU:         getString(productData, "sku"),
		Name:        getString(productData, "name"),
		BrandName:   getString(productData, "brand_name"),
		Gender:      getString(productData, "gender"),
		Category:    getString(productData, "category"),
		Color:       getString(productData, "color"),
		Size:        getString(productData, "size"),
		Material:    getString(productData, "material"),
		Image1:      getString(productData, "image1"),
		Image2:      getString(productData, "image2"),
		Description: getString(productData, "description"),
		Data:        models.JSONMap(productData),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Handle numeric fields
	if price := getFloat(productData, "price"); price != nil {
		product.Price = *price
	}
	if mrp := getFloat(productData, "mrp"); mrp != nil {
		product.MRP = *mrp
	}
	if qty := getInt(productData, "quantity"); qty != nil {
		product.Quantity = int(*qty)
	}

	// Validate required fields
	if product.SKU == "" {
		return nil, errors.New("SKU is required")
	}
	if product.Name == "" {
		return nil, errors.New("Name is required")
	}

	return product, nil
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok && val != nil {
		str := fmt.Sprintf("%v", val)
		return strings.TrimSpace(str)
	}
	return ""
}

func getFloat(data map[string]interface{}, key string) *float64 {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case float64:
			return &v
		case string:
			v = strings.TrimSpace(v)
			if v == "" {
				return nil
			}
			var f float64
			if _, err := fmt.Sscanf(v, "%f", &f); err == nil {
				return &f
			}
		case int:
			f := float64(v)
			return &f
		}
	}
	return nil
}

func getInt(data map[string]interface{}, key string) *int64 {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case int64:
			return &v
		case float64:
			i := int64(v)
			return &i
		case string:
			v = strings.TrimSpace(v)
			if v == "" {
				return nil
			}
			// Check if string contains decimal point - if yes, it's not a valid int
			if strings.Contains(v, ".") {
				return nil
			}
			var i int64
			if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
				return &i
			}
		case int:
			i := int64(v)
			return &i
		}
	}
	return nil
}
