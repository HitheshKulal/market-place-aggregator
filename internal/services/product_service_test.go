package services

import (
	"bytes"
	"mime/multipart"
	"testing"

	"go-backend-api/internal/models"

	"github.com/stretchr/testify/assert"
)

// ==================== EXTRACTION TESTS ====================

func TestExtractCSVData(t *testing.T) {
	tests := []struct {
		name             string
		csvContent       string
		expectedHeaders  []string
		expectedRowCount int
		expectedError    bool
		errorMessage     string
	}{
		{
			name: "Valid CSV with multiple rows",
			csvContent: `sku,name,brand_name,price,quantity
IPHONE-001,iPhone 15 Pro,Apple,999.99,50
SAMSUNG-001,Galaxy S24,Samsung,899.99,75
SONY-001,WH-1000XM5,Sony,349.99,120`,
			expectedHeaders:  []string{"sku", "name", "brand_name", "price", "quantity"},
			expectedRowCount: 3,
			expectedError:    false,
		},
		{
			name: "CSV with whitespace in headers",
			csvContent: `  sku  , name ,brand_name,  price  
PROD-001,Product One,Brand A,99.99`,
			expectedHeaders:  []string{"sku", "name", "brand_name", "price"},
			expectedRowCount: 1,
			expectedError:    false,
		},
		{
			name:          "Empty CSV file",
			csvContent:    "",
			expectedError: true,
			errorMessage:  "CSV file is empty",
		},
		{
			name:          "CSV with only headers (no data)",
			csvContent:    "sku,name,price",
			expectedError: true,
			errorMessage:  "CSV file must have at least headers and one data row",
		},
		{
			name: "CSV with empty column header",
			csvContent: `sku,,price
PROD-001,Product,99.99`,
			expectedError: true,
			errorMessage:  "CSV contains empty column headers",
		},
		{
			name: "CSV with quoted values",
			csvContent: `sku,name,description
PROD-001,"Product Name","Description with, comma"
PROD-002,Simple Product,Simple description`,
			expectedHeaders:  []string{"sku", "name", "description"},
			expectedRowCount: 2,
			expectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := createMockFile(tt.csvContent)
			service := &ProductService{}

			headers, dataRows, err := service.extractCSVData(file)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedHeaders, headers)
				assert.Equal(t, tt.expectedRowCount, len(dataRows))
			}
		})
	}
}

// ==================== MAPPING TESTS ====================

func TestMapRowToProductData(t *testing.T) {
	tests := []struct {
		name     string
		headers  []string
		row      []string
		expected map[string]interface{}
	}{
		{
			name:    "Complete row mapping",
			headers: []string{"sku", "name", "brand_name", "price", "quantity"},
			row:     []string{"IPHONE-001", "iPhone 15 Pro", "Apple", "999.99", "50"},
			expected: map[string]interface{}{
				"sku":        "IPHONE-001",
				"name":       "iPhone 15 Pro",
				"brand_name": "Apple",
				"price":      "999.99",
				"quantity":   "50",
			},
		},
		{
			name:    "Row with whitespace values",
			headers: []string{"sku", "name", "price"},
			row:     []string{"  PROD-001  ", "  Product Name  ", "  99.99  "},
			expected: map[string]interface{}{
				"sku":   "PROD-001",
				"name":  "Product Name",
				"price": "99.99",
			},
		},
		{
			name:    "Row with empty values",
			headers: []string{"sku", "name", "brand_name", "price"},
			row:     []string{"PROD-002", "Product Two", "", "199.99"},
			expected: map[string]interface{}{
				"sku":        "PROD-002",
				"name":       "Product Two",
				"brand_name": "",
				"price":      "199.99",
			},
		},
		{
			name:    "Row with fewer values than headers",
			headers: []string{"sku", "name", "brand_name", "price", "quantity"},
			row:     []string{"PROD-003", "Product Three", "Brand"},
			expected: map[string]interface{}{
				"sku":        "PROD-003",
				"name":       "Product Three",
				"brand_name": "Brand",
			},
		},
		{
			name:    "Row with more values than headers (extra values ignored)",
			headers: []string{"sku", "name", "price"},
			row:     []string{"PROD-004", "Product Four", "49.99", "ExtraValue1", "ExtraValue2"},
			expected: map[string]interface{}{
				"sku":   "PROD-004",
				"name":  "Product Four",
				"price": "49.99",
			},
		},
		{
			name:    "Empty row",
			headers: []string{"sku", "name", "price"},
			row:     []string{"", "", ""},
			expected: map[string]interface{}{
				"sku":   "",
				"name":  "",
				"price": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ProductService{}
			result := service.mapRowToProductData(tt.headers, tt.row)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ==================== VALIDATION TESTS ====================

func TestValidateAndCreateProduct(t *testing.T) {
	tests := []struct {
		name          string
		productData   map[string]interface{}
		expectedError bool
		errorMessage  string
		validateFunc  func(*testing.T, *models.Product)
	}{
		{
			name: "Valid product with all fields",
			productData: map[string]interface{}{
				"sku":         "IPHONE-001",
				"name":        "iPhone 15 Pro",
				"brand_name":  "Apple",
				"gender":      "Unisex",
				"category":    "Electronics",
				"color":       "Titanium",
				"size":        "256GB",
				"mrp":         "1199.99",
				"price":       "999.99",
				"material":    "Titanium",
				"image1":      "https://example.com/image1.jpg",
				"image2":      "https://example.com/image2.jpg",
				"quantity":    "50",
				"description": "Latest iPhone",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, p *models.Product) {
				assert.Equal(t, "IPHONE-001", p.SKU)
				assert.Equal(t, "iPhone 15 Pro", p.Name)
				assert.Equal(t, "Apple", p.BrandName)
				assert.Equal(t, "Unisex", p.Gender)
				assert.Equal(t, "Electronics", p.Category)
				assert.Equal(t, 999.99, p.Price)
				assert.Equal(t, 1199.99, p.MRP)
				assert.Equal(t, 50, p.Quantity)
			},
		},
		{
			name: "Missing required SKU field",
			productData: map[string]interface{}{
				"name":  "Product Without SKU",
				"price": "99.99",
			},
			expectedError: true,
			errorMessage:  "SKU is required",
		},
		{
			name: "Missing required Name field",
			productData: map[string]interface{}{
				"sku":   "PROD-001",
				"price": "99.99",
			},
			expectedError: true,
			errorMessage:  "Name is required",
		},
		{
			name: "Empty SKU",
			productData: map[string]interface{}{
				"sku":  "",
				"name": "Product",
			},
			expectedError: true,
			errorMessage:  "SKU is required",
		},
		{
			name: "Empty Name",
			productData: map[string]interface{}{
				"sku":  "PROD-002",
				"name": "",
			},
			expectedError: true,
			errorMessage:  "Name is required",
		},
		{
			name: "Valid product with minimal fields",
			productData: map[string]interface{}{
				"sku":  "PROD-003",
				"name": "Minimal Product",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, p *models.Product) {
				assert.Equal(t, "PROD-003", p.SKU)
				assert.Equal(t, "Minimal Product", p.Name)
				assert.Equal(t, 0.0, p.Price)
				assert.Equal(t, 0.0, p.MRP)
				assert.Equal(t, 0, p.Quantity)
			},
		},
		{
			name: "Numeric fields as strings",
			productData: map[string]interface{}{
				"sku":      "PROD-004",
				"name":     "Product Four",
				"price":    "199.99",
				"mrp":      "299.99",
				"quantity": "100",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, p *models.Product) {
				assert.Equal(t, 199.99, p.Price)
				assert.Equal(t, 299.99, p.MRP)
				assert.Equal(t, 100, p.Quantity)
			},
		},
		{
			name: "Numeric fields as actual numbers",
			productData: map[string]interface{}{
				"sku":      "PROD-005",
				"name":     "Product Five",
				"price":    149.99,
				"mrp":      249.99,
				"quantity": 75,
			},
			expectedError: false,
			validateFunc: func(t *testing.T, p *models.Product) {
				assert.Equal(t, 149.99, p.Price)
				assert.Equal(t, 249.99, p.MRP)
				assert.Equal(t, 75, p.Quantity)
			},
		},
		{
			name: "Empty numeric fields",
			productData: map[string]interface{}{
				"sku":      "PROD-006",
				"name":     "Product Six",
				"price":    "",
				"mrp":      "",
				"quantity": "",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, p *models.Product) {
				assert.Equal(t, 0.0, p.Price)
				assert.Equal(t, 0.0, p.MRP)
				assert.Equal(t, 0, p.Quantity)
			},
		},
		{
			name: "Whitespace trimming in fields",
			productData: map[string]interface{}{
				"sku":        "  PROD-007  ",
				"name":       "  Product Seven  ",
				"brand_name": "  Brand Name  ",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, p *models.Product) {
				assert.Equal(t, "PROD-007", p.SKU)
				assert.Equal(t, "Product Seven", p.Name)
				assert.Equal(t, "Brand Name", p.BrandName)
			},
		},
		{
			name: "Product with extra fields in Data",
			productData: map[string]interface{}{
				"sku":          "PROD-008",
				"name":         "Product Eight",
				"custom_field": "Custom Value",
				"extra_data":   "Extra",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, p *models.Product) {
				assert.Equal(t, "PROD-008", p.SKU)
				assert.Equal(t, "Product Eight", p.Name)
				assert.NotNil(t, p.Data)
				assert.Equal(t, "Custom Value", p.Data["custom_field"])
				assert.Equal(t, "Extra", p.Data["extra_data"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ProductService{}
			product, err := service.validateAndCreateProduct(tt.productData)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				assert.Nil(t, product)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, product)
				if tt.validateFunc != nil {
					tt.validateFunc(t, product)
				}
			}
		})
	}
}

// ==================== HELPER FUNCTION TESTS ====================

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected string
	}{
		{"String value", map[string]interface{}{"key": "value"}, "key", "value"},
		{"Integer value", map[string]interface{}{"key": 123}, "key", "123"},
		{"Float value", map[string]interface{}{"key": 45.67}, "key", "45.67"},
		{"Empty string", map[string]interface{}{"key": ""}, "key", ""},
		{"Missing key", map[string]interface{}{}, "key", ""},
		{"Nil value", map[string]interface{}{"key": nil}, "key", ""},
		{"Whitespace string", map[string]interface{}{"key": "  value  "}, "key", "value"},
		{"Boolean true", map[string]interface{}{"key": true}, "key", "true"},
		{"Boolean false", map[string]interface{}{"key": false}, "key", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getString(tt.data, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetFloat(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected *float64
	}{
		{"Float64 value", map[string]interface{}{"key": 99.99}, "key", floatPtr(99.99)},
		{"String float", map[string]interface{}{"key": "99.99"}, "key", floatPtr(99.99)},
		{"Integer value", map[string]interface{}{"key": 100}, "key", floatPtr(100.0)},
		{"String integer", map[string]interface{}{"key": "100"}, "key", floatPtr(100.0)},
		{"Empty string", map[string]interface{}{"key": ""}, "key", nil},
		{"Invalid string", map[string]interface{}{"key": "invalid"}, "key", nil},
		{"Missing key", map[string]interface{}{}, "key", nil},
		{"Nil value", map[string]interface{}{"key": nil}, "key", nil},
		{"Whitespace string", map[string]interface{}{"key": "  123.45  "}, "key", floatPtr(123.45)},
		{"Negative float", map[string]interface{}{"key": "-99.99"}, "key", floatPtr(-99.99)},
		{"Zero", map[string]interface{}{"key": "0"}, "key", floatPtr(0.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFloat(tt.data, tt.key)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.InDelta(t, *tt.expected, *result, 0.001)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected *int64
	}{
		{"Int64 value", map[string]interface{}{"key": int64(100)}, "key", intPtr(100)},
		{"Float64 value", map[string]interface{}{"key": float64(100.9)}, "key", intPtr(100)},
		{"String int", map[string]interface{}{"key": "100"}, "key", intPtr(100)},
		{"Regular int", map[string]interface{}{"key": 50}, "key", intPtr(50)},
		{"Empty string", map[string]interface{}{"key": ""}, "key", nil},
		{"Invalid string", map[string]interface{}{"key": "invalid"}, "key", nil},
		{"Missing key", map[string]interface{}{}, "key", nil},
		{"Nil value", map[string]interface{}{"key": nil}, "key", nil},
		{"Whitespace string", map[string]interface{}{"key": "  75  "}, "key", intPtr(75)},
		{"Negative int", map[string]interface{}{"key": "-50"}, "key", intPtr(-50)},
		{"Zero", map[string]interface{}{"key": "0"}, "key", intPtr(0)},
		{"Float string (should fail)", map[string]interface{}{"key": "99.99"}, "key", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getInt(tt.data, tt.key)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

// ==================== INTEGRATION TEST ====================

func TestFullCSVProcessingFlow(t *testing.T) {
	csvContent := `sku,name,brand_name,price,quantity
IPHONE-001,iPhone 15 Pro,Apple,999.99,50
SAMSUNG-001,Galaxy S24,Samsung,899.99,75
INVALID-001,,Sony,349.99,120`

	file := createMockFile(csvContent)
	service := &ProductService{}

	// Step 1: Extract
	headers, dataRows, err := service.extractCSVData(file)
	assert.NoError(t, err)
	assert.Equal(t, []string{"sku", "name", "brand_name", "price", "quantity"}, headers)
	assert.Equal(t, 3, len(dataRows))

	// Step 2: Map and Validate
	validProducts := 0
	invalidProducts := 0

	for _, row := range dataRows {
		productData := service.mapRowToProductData(headers, row)
		product, err := service.validateAndCreateProduct(productData)

		if err != nil {
			invalidProducts++
		} else {
			validProducts++
			assert.NotNil(t, product)
		}
	}

	assert.Equal(t, 2, validProducts)
	assert.Equal(t, 1, invalidProducts)
}

// ==================== HELPER FUNCTIONS ====================

func createMockFile(content string) multipart.File {
	return &mockFile{
		Reader: bytes.NewReader([]byte(content)),
	}
}

type mockFile struct {
	*bytes.Reader
}

func (m *mockFile) Close() error {
	return nil
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int64) *int64 {
	return &i
}
