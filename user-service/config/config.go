package config

import (
	"log"

	"github.com/faris-muhammed/e-commerce/user-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase connects to the PostgreSQL database for both services
func ConnectDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=501417 dbname=user_service port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = database
	// Auto migrate models for both Category and Seller
	if err := DB.AutoMigrate(&models.User{},&models.OTP{},models.Cart{},models.Address{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully!")
	return DB, nil
}
