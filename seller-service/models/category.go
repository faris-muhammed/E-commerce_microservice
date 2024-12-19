package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID          int32 `gorm:"primary_key"`
	Name        string
	Description string
}
