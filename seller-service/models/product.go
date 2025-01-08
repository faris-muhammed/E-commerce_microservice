package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string
	Price      float64
	Stock      uint64
	IsDeleted  bool
	CategoryID uint64
	SellerID   uint64
}
