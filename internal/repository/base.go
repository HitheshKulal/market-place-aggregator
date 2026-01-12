package repository

import (
	"context"

	"fmt"

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

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *Repository[T]) GetByIDs(ctx context.Context, ids []uint) ([]T, error) {
	var entities []T
	if len(ids) == 0 {
		return entities, nil
	}

	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *Repository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T

	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Repository[T]) DeleteByID(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

func (r *Repository[T]) RawQuery(ctx context.Context, query string, args ...interface{}) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) ApplyOrderByFilter(query *gorm.DB, orderBy string, orderDirection string) *gorm.DB {
	if orderBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDirection))
	}
	return query
}

func (r *Repository[T]) ApplyPaginationFilter(query *gorm.DB, page int, limit int) *gorm.DB {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	return query
}

func (r *Repository[T]) GetDB() *gorm.DB {
	return r.db
}

// Repository registry
type Repositories struct {
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{}
}
