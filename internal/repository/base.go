package repository

import (
	"gorm.io/gorm"
)

// Base repository for common CRUD operations
type Repository[T any] struct {
	db *gorm.DB
}

type QueryOptions struct {
	OrderBy        string
	OrderDirection string
	Page           int
	Limit          int
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

type Repositories struct {
	Product  *ProductRepository
	Template *TemplateRepository
	Mapping  *MappingRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Product:  NewProductRepository(db),
		Template: NewTemplateRepository(db),
		Mapping:  NewMappingRepository(db),
	}
}
