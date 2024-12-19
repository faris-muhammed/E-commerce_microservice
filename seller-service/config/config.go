package config

import (
	"log"
	"net"

	"github.com/faris-muhammed/e-commerce/seller-service/handlers"
	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"github.com/faris-muhammed/e-commerce/seller-service/repository"
	"github.com/faris-muhammed/e-commerce/seller-service/service"
	sellerpb "github.com/faris-muhammed/e-protofiles/seller"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=501417 dbname=seller_service port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = database
	if err := DB.AutoMigrate(&models.Product{}, &models.Category{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully!")
	return DB, nil
}

// setupDependencies initializes all required dependencies.
func SetupDependencies() (*handlers.CategoryHandler, error) {
	// Step 1: Connect to the database
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}

	categoryRepo := repository.NewSellerRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	return handlers.NewCategoryHandler(categoryService), nil

}

func StartGRPCServer(handler sellerpb.CategoryServiceServer) {
	// Step 1: Create gRPC server
	grpcServer := grpc.NewServer()
	sellerpb.RegisterCategoryServiceServer(grpcServer, handler)

	// Step 2: Listen on port 50052
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Seller service is running on port 50052...")

	// Step 3: Serve gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
