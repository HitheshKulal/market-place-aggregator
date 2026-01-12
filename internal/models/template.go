package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;unique" json:"name"`
	Fields    StringArray    `gorm:"type:jsonb;not null" json:"fields"` // ["title", "price", "quantity"]
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// StringArray is a custom type for []string that works with JSONB
type StringArray []string

// Scan implements sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value for StringArray")
	}

	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(s)
}

func AutoMigrateMapping(db *gorm.DB) error {
	return db.AutoMigrate(&Mapping{})
}
