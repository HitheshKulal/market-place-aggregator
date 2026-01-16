package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Mapping struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TemplateID uint           `gorm:"index;not null;" json:"templateId"`   // One mapping per template
	SellerID   uint           `gorm:"index;not null;" json:"sellerId"`     // One mapping per template
	FieldMap   StringMap      `gorm:"type:jsonb;not null" json:"fieldMap"` // {"title":"name", "price":"price", "quantity":"quantity"}
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Template *Template `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
	Seller   *Seller   `gorm:"foreignKey:SellerID;constraint:OnDelete:CASCADE" json:"seller,omitempty"`
}

// StringMap is a custom type for map[string]string that works with JSONB
type StringMap map[string]string

// Scan implements sql.Scanner interface
func (m *StringMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(map[string]string)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value for StringMap")
	}

	return json.Unmarshal(bytes, m)
}

// Value implements driver.Valuer interface
func (m StringMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return json.Marshal(map[string]string{})
	}

	return json.Marshal(m)
}

func AutoMigrateTemplate(db *gorm.DB) error {
	return db.AutoMigrate(&Template{})
}
