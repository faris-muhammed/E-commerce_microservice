package models

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	Username string
	Password string
}