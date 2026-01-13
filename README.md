# Product Template Mapping System

A robust Go-based API system for managing products with dynamic template mappings. This system allows you to define templates (like Amazon, Flipkart formats), map product fields to template fields, and transform product data according to different marketplace requirements.

## Table of Contents

- [System Design](#system-design)
- [Database Schema](#database-schema)
- [API Documentation](#api-documentation)
- [Setup Instructions](#setup-instructions)
- [Running with Docker](#running-with-docker)
- [Running Tests](#running-tests)
- [Usage Examples](#usage-examples)

---

## System Design

### Architecture Overview
```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│         API Layer (Gin)             │
│  ┌──────────┐  ┌──────────┐       │
│  │ Handlers │  │  Routes  │       │
│  └──────────┘  └──────────┘       │
└──────────┬──────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│       Service Layer                 │
│  ┌──────────────┐  ┌─────────────┐ │
│  │   Product    │  │  Template   │ │
│  │   Service    │  │   Service   │ │
│  └──────────────┘  └─────────────┘ │
│  ┌──────────────┐                  │
│  │   Mapping    │                  │
│  │   Service    │                  │
│  └──────────────┘                  │
└──────────┬──────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│      Repository Layer               │
│  ┌──────────────┐  ┌─────────────┐ │
│  │   Product    │  │  Template   │ │
│  │   Repo       │  │   Repo      │ │
│  └──────────────┘  └─────────────┘ │
│  ┌──────────────┐                  │
│  │   Mapping    │                  │
│  │   Repo       │                  │
│  └──────────────┘                  │
└──────────┬──────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│       PostgreSQL Database           │
│  ┌──────────┐  ┌──────────┐        │
│  │ Products │  │Templates │        │
│  └──────────┘  └──────────┘        │
│  ┌──────────┐                      │
│  │ Mappings │                      │
│  └──────────┘                      │
└─────────────────────────────────────┘
```

### Key Components

1. **Products**: Store product information with flexible schema
2. **Templates**: Define output formats for different marketplaces
3. **Mappings**: Define field transformations from products to templates
4. **Transformation Engine**: Converts product data according to template mappings

### Data Flow
```
CSV Upload → Parse → Validate → Map → Store → Transform → Output
```

---
## Database Schema

### Products Table
```sql
CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    sku         VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    brand_name  VARCHAR(255),
    gender      VARCHAR(50),
    category    VARCHAR(100),
    color       VARCHAR(50),
    size        VARCHAR(50),
    mrp         NUMERIC,
    price       NUMERIC,
    material    VARCHAR(100),
    image1      VARCHAR(255),
    image2      VARCHAR(255),
    quantity    BIGINT,
    description TEXT,
    data        JSONB NOT NULL,           -- Stores all fields including extras
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE UNIQUE INDEX uni_products_sku ON products(sku);
```

### Templates Table
```sql
CREATE TABLE templates (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    fields     JSONB NOT NULL,            -- Array of field names
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_templates_deleted_at ON templates(deleted_at);
CREATE UNIQUE INDEX uni_templates_name ON templates(name);
```

### Mappings Table
```sql
CREATE TABLE mappings (
    id          BIGSERIAL PRIMARY KEY,
    template_id BIGINT NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    field_map   JSONB NOT NULL,           -- Field mapping object
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_mappings_deleted_at ON mappings(deleted_at);
CREATE INDEX idx_mappings_template_id ON mappings(template_id);
-- Note: No unique constraint on template_id (allows 1:N relationship)
```

### Entity Relationship Diagram
```
┌─────────────────┐
│    Templates    │
│─────────────────│
│ id (PK)         │
│ name (UNIQUE)   │
│ fields (JSONB)  │
└────────┬────────┘
         │
         │ 1:N
         │ (One template can have multiple mappings)
         │
         ▼
┌─────────────────┐
│    Mappings     │
│─────────────────│
│ id (PK)         │
│ template_id(FK) │
│ field_map(JSONB)│
└─────────────────┘

┌─────────────────┐
│    Products     │
│─────────────────│
│ id (PK)         │
│ sku (UNIQUE)    │
│ name            │
│ brand_name      │
│ gender          │
│ category        │
│ color           │
│ size            │
│ mrp             │
│ price           │
│ material        │
│ image1          │
│ image2          │
│ quantity        │
│ description     │
│ data (JSONB)    │
└─────────────────┘

Relationship: Templates (1) ──< (N) Mappings
One template can have multiple mappings for different use cases
```

### Relationship Details

- **Templates ↔ Mappings**: One-to-Many (1:N)
  - One template (e.g., "Flipkart Template") can have multiple mappings
  - Each mapping references one template via `template_id` foreign key
  - Cascade delete: When a template is deleted, all its mappings are deleted
  
- **Products**: Independent table
  - No direct foreign key relationship with templates or mappings
  - Products are transformed at runtime based on selected mapping

### Example Data Structure
```sql
-- One Template
INSERT INTO templates (id, name, fields) VALUES 
(1, 'Flipkart Template', '["productName", "brand", "price", "sku"]');

-- Multiple Mappings for the same template
INSERT INTO mappings (template_id, field_map) VALUES 
(1, '{"productName": "name", "brand": "brand_name", "price": "price", "sku": "sku"}'),
(1, '{"productName": "name", "brand": "brand_name", "price": "mrp", "sku": "sku"}'),
(1, '{"productName": "name", "brand": "brand_name", "price": "price", "sku": "sku"}');
```
## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### 1. Products

##### Get All Products
```http
GET /api/v1/products/
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "sku": "IPHONE-15-001",
      "name": "iPhone 15 Pro",
      "brandName": "Apple",
      "gender": "Unisex",
      "category": "Electronics",
      "color": "Titanium",
      "size": "256GB",
      "mrp": 1199.99,
      "price": 999.99,
      "material": "Titanium",
      "image1": "https://example.com/image1.jpg",
      "image2": "https://example.com/image2.jpg",
      "quantity": 50,
      "description": "Latest iPhone",
      "data": {},
      "createdAt": "2024-01-12T10:00:00Z",
      "updatedAt": "2024-01-12T10:00:00Z"
    }
  ],
  "statusCode": 200
}
```

##### Upload Products (CSV)
```http
POST /api/v1/products/upload
Content-Type: multipart/form-data
```

**Request:**
```
file: products.csv (form-data)
```

**Sample CSV Format:**
```csv
sku,name,brand_name,price,quantity,category
IPHONE-001,iPhone 15 Pro,Apple,999.99,50,Electronics
SAMSUNG-001,Galaxy S24,Samsung,899.99,75,Electronics
```

**Response:**
```json
{
  "success": true,
  "message": "Products uploaded and stored",
  "data": {
    "fileName": "products.csv",
    "totalRows": 10,
    "successCount": 9,
    "failedCount": 1,
    "discoveredColumns": [
      "sku",
      "name",
      "brand_name",
      "price",
      "quantity"
    ],
    "sampleProducts": [
      {
        "id": 1,
        "sku": "IPHONE-001",
        "name": "iPhone 15 Pro",
        "price": 999.99,
        "quantity": 50
      }
    ],
    "errors": [
      {
        "row": 5,
        "sku": "DUPLICATE-SKU",
        "message": "Database error: duplicate key value violates unique constraint"
      }
    ]
  },
  "statusCode": 201
}
```

**Status Codes:**
- `201`: Products created successfully
- `206`: Partial success (some products failed)
- `400`: Invalid request
- `422`: Validation error

#### 2. Templates

##### Get All Templates
```http
GET /api/v1/templates/
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Flipkart Template",
      "fields": [
        "productName",
        "brand",
        "gender",
        "category",
        "color",
        "size",
        "mrp",
        "price",
        "sku",
        "description",
        "material",
        "images"
      ],
      "createdAt": "2024-01-12T10:00:00Z",
      "updatedAt": "2024-01-12T10:00:00Z"
    }
  ],
  "statusCode": 200
}
```

##### Create Template
```http
POST /api/v1/templates/
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Amazon Template",
  "fields": [
    "title",
    "brand",
    "price",
    "quantity",
    "description"
  ],
  "isActive": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 2,
    "name": "Amazon Template",
    "fields": [
      "title",
      "brand",
      "price",
      "quantity",
      "description"
    ],
    "createdAt": "2024-01-12T10:30:00Z",
    "updatedAt": "2024-01-12T10:30:00Z"
  },
  "statusCode": 200
}
```

#### 3. Mappings

##### Get All Mappings
```http
GET /api/v1/mappings/
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "templateId": 1,
      "fieldMap": {
        "productName": "name",
        "brand": "brand_name",
        "gender": "gender",
        "category": "category",
        "color": "color",
        "size": "size",
        "mrp": "mrp",
        "price": "price",
        "sku": "sku",
        "description": "description",
        "material": "material",
        "images": "image1"
      },
      "createdAt": "2024-01-12T10:00:00Z",
      "updatedAt": "2024-01-12T10:00:00Z",
      "template": {
        "id": 1,
        "name": "Flipkart Template",
        "fields": [
          "productName",
          "brand",
          "price"
        ]
      }
    }
  ],
  "statusCode": 200
}
```

##### Create Mapping
```http
POST /api/v1/mappings/
Content-Type: application/json
```

**Request Body:**
```json
{
  "templateId": 1,
  "fieldMap": {
    "productName": "name",
    "brand": "brand_name",
    "price": "price",
    "quantity": "quantity"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 2,
    "templateId": 1,
    "fieldMap": {
      "productName": "name",
      "brand": "brand_name",
      "price": "price",
      "quantity": "quantity"
    },
    "createdAt": "2024-01-12T10:45:00Z",
    "updatedAt": "2024-01-12T10:45:00Z"
  },
  "statusCode": 200
}
```

##### Get Products by Template (Transformed)
```http
GET /api/v1/mappings/:id/products
```

**Example:**
```http
GET /api/v1/mappings/1/products
```

**Response (Products transformed according to Flipkart template):**
```json
{
  "success": true,
  "data": [
    {
      "productName": "iPhone 15 Pro",
      "brand": "Apple",
      "gender": "Unisex",
      "category": "Electronics",
      "color": "Titanium",
      "size": "256GB",
      "mrp": 1199.99,
      "price": 999.99,
      "sku": "IPHONE-15-001",
      "description": "Latest iPhone",
      "material": "Titanium",
      "images": "https://example.com/image1.jpg"
    },
    {
      "productName": "Galaxy S24",
      "brand": "Samsung",
      "gender": "Unisex",
      "category": "Electronics",
      "color": "Black",
      "size": "128GB",
      "mrp": 999.99,
      "price": 899.99,
      "sku": "SAMSUNG-001",
      "description": "Galaxy AI smartphone",
      "material": "Glass",
      "images": "https://example.com/image2.jpg"
    }
  ],
  "statusCode": 200
}
```

---

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker & Docker Compose (for containerized setup)

### Local Development Setup

#### 1. Clone the repository
```bash
git clone <repository-url>
cd go-backend-api
```

#### 2. Create environment file
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your configurations
nano .env
```

**Required Environment Variables:**
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=product_db
DB_SSLMODE=disable

# Server
PORT=8080
GIN_MODE=debug

# Application
APP_ENV=development
```

#### 3. Install dependencies
```bash
go mod download
```

#### 4. Setup PostgreSQL Database
```bash
# Create database
psql -U postgres -c "CREATE DATABASE product_db;"

# Run migrations (auto-migration will run on app start)
```

#### 5. Run the application
```bash
go run cmd/server/main.go
```

The server will start at `http://localhost:8080`

---

## Running with Docker

### Option 1: Docker Compose (Recommended)

#### 1. Create environment file
```bash
# Copy example environment file
cp .env.example .env
```

**Update .env for Docker:**
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=product_db
DB_SSLMODE=disable

PORT=8080
GIN_MODE=release
APP_ENV=production
```

#### 2. Build and run with Docker Compose
```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

#### 3. Access the application
```
API: http://localhost:8080
Database: localhost:5432
```

#### Run tests with coverage
```bash
# Generate coverage report
go test ./internal/services -coverprofile=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open coverage report in browser
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

#### Run tests in Docker
```bash
# Run tests in container
docker-compose run --rm app go test ./internal/services -v

# Run with coverage
docker-compose run --rm app go test ./internal/services -cover
```

### Test Coverage Summary
```
Extraction Tests:     6 test cases
Mapping Tests:        6 test cases
Validation Tests:    11 test cases
Helper Functions:    28 test cases
Integration Tests:    1 test case
─────────────────────────────────
Total:               52 test cases
Coverage:            ~42% of statements
```

### Load Testing (Optional)
```bash
# Install Apache Bench
sudo apt-get install apache2-utils

# Test product upload endpoint
ab -n 100 -c 10 -p products.csv -T 'multipart/form-data' \
  http://localhost:8080/api/v1/products/upload

# Test get products endpoint
ab -n 1000 -c 50 http://localhost:8080/api/v1/products/
```

---

## Usage Examples

### Complete Workflow Example

#### Step 1: Create a Template
```bash
curl -X POST http://localhost:8080/api/v1/templates/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Flipkart Template",
    "fields": ["productName", "brand", "price", "sku"]
  }'
```

#### Step 2: Create a Mapping
```bash
curl -X POST http://localhost:8080/api/v1/mappings/ \
  -H "Content-Type: application/json" \
  -d '{
    "templateId": 1,
    "fieldMap": {
      "productName": "name",
      "brand": "brand_name",
      "price": "price",
      "sku": "sku"
    }
  }'
```

#### Step 3: Upload Products
```bash
curl -X POST http://localhost:8080/api/v1/products/upload \
  -F "file=@products.csv"
```

#### Step 4: Get Transformed Products
```bash
curl -X GET http://localhost:8080/api/v1/mappings/1/products
```

### Sample CSV File

Create `products.csv`:
```csv
sku,name,brand_name,price,quantity,category
IPHONE-001,iPhone 15 Pro,Apple,999.99,50,Electronics
SAMSUNG-001,Galaxy S24,Samsung,899.99,75,Electronics
SONY-001,WH-1000XM5,Sony,349.99,120,Audio
```

### Postman Collection

Import the following JSON into Postman:
```json
{
  "info": {
    "name": "Product Template API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Products",
      "item": [
        {
          "name": "Get All Products",
          "request": {
            "method": "GET",
            "url": "{{base_url}}/products/"
          }
        },
        {
          "name": "Upload Products CSV",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/products/upload",
            "body": {
              "mode": "formdata",
              "formdata": [
                {
                  "key": "file",
                  "type": "file",
                  "src": "/path/to/products.csv"
                }
              ]
            }
          }
        }
      ]
    },
    {
      "name": "Templates",
      "item": [
        {
          "name": "Get All Templates",
          "request": {
            "method": "GET",
            "url": "{{base_url}}/templates/"
          }
        },
        {
          "name": "Create Template",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/templates/",
            "header": [
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"name\": \"Flipkart Template\",\n  \"fields\": [\"productName\", \"brand\", \"price\"]\n}"
            }
          }
        }
      ]
    },
    {
      "name": "Mappings",
      "item": [
        {
          "name": "Get All Mappings",
          "request": {
            "method": "GET",
            "url": "{{base_url}}/mappings/"
          }
        },
        {
          "name": "Create Mapping",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/mappings/",
            "header": [
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"templateId\": 1,\n  \"fieldMap\": {\n    \"productName\": \"name\",\n    \"brand\": \"brand_name\",\n    \"price\": \"price\"\n  }\n}"
            }
          }
        },
        {
          "name": "Get Products by Template",
          "request": {
            "method": "GET",
            "url": "{{base_url}}/mappings/1/products"
          }
        }
      ]
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080/api/v1"
    }
  ]
}
```

---

## Troubleshooting

### Common Issues

#### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs postgres

# Verify connection
psql -h localhost -U postgres -d product_db
```

#### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port in .env
PORT=8081
```

#### CSV Upload Fails

- Ensure CSV has headers in the first row
- Check that SKU and Name columns are present
- Verify file encoding is UTF-8
- Check file size limits

---

## Project Structure
```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── product/
│   │   ├── template/
│   │   └── mapping/
│   ├── services/                # Business logic
│   │   ├── product_service.go
│   │   ├── template_service.go
│   │   └── mapping_service.go
│   ├── repository/              # Data access layer
│   │   ├── product_repository.go
│   │   ├── template_repository.go
│   │   └── mapping_repository.go
│   ├── models/                  # Database models
│   │   ├── product.go
│   │   ├── template.go
│   │   └── mapping.go
│   └── requests/                # Request/Response DTOs
├── .env.example                 # Example environment file
├── docker-compose.yml           # Docker compose configuration
├── Dockerfile                   # Docker image definition
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
└── README.md                    # This file
```