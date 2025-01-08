package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID    uint64 `gorm:"not null"`
	ProductID uint64 `gorm:"not null"`
	Quantity  uint32 `gorm:"not null"`
}
