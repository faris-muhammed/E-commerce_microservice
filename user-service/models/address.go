package models

type Address struct {
	ID         uint64 `gorm:"primaryKey"`
	UserID     uint64 `gorm:"not null"`
	Street     string `gorm:"not null"`
	City       string `gorm:"not null"`
	State      string `gorm:"not null"`
	PostalCode string `gorm:"not null"`
	Country    string `gorm:"not null"`
}
