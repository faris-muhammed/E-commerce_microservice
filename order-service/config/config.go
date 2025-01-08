package config

import (
	"log"

	"github.com/faris-muhammed/e-commerce/order-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=501417 dbname=order_service port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = database
	// Auto migrate models for both Category and Seller
	if err := DB.AutoMigrate(&models.Order{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully!")
	return DB, nil
}
