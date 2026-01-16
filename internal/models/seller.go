package models

import "gorm.io/gorm"

type Seller struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

func AutoMigrateSeller(db *gorm.DB) error {
	return db.AutoMigrate(&Seller{})
}
