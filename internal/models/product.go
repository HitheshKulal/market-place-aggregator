package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SKU         string         `gorm:"size:100;not null;unique" json:"sku"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	BrandName   string         `gorm:"size:255" json:"brandName"`
	Gender      string         `gorm:"size:50" json:"gender"`
	Category    string         `gorm:"size:100" json:"category"`
	Color       string         `gorm:"size:50" json:"color"`
	Size        string         `gorm:"size:50" json:"size"`
	MRP         float64        `json:"mrp"`
	Price       float64        `json:"price"`
	Material    string         `gorm:"size:100" json:"material"`
	Image1      string         `gorm:"size:255" json:"image1"`
	Image2      string         `gorm:"size:255" json:"image2"`
	Quantity    int            `json:"quantity"`
	Description string         `gorm:"type:text" json:"description"`
	Data        JSONMap        `gorm:"type:jsonb;not null" json:"data"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// JSONMap is a custom type for map[string]interface{} that works with JSONB
type JSONMap map[string]interface{}

// Scan implements sql.Scanner interface
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value for JSONMap")
	}

	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return json.Marshal(map[string]interface{}{})
	}

	return json.Marshal(j)
}

func AutoMigrateProduct(db *gorm.DB) error {
	return db.AutoMigrate(&Product{})
}
