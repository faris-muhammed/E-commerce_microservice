package config

import (
	"log"
	"net"

	"github.com/faris-muhammed/e-commerce/seller-service/handlers"
	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"github.com/faris-muhammed/e-commerce/seller-service/repository"
	"github.com/faris-muhammed/e-commerce/seller-service/service"
	categorypb "github.com/faris-muhammed/e-protofiles/category"
	productpb "github.com/faris-muhammed/e-protofiles/product"
	sellerpb "github.com/faris-muhammed/e-protofiles/sellerlogin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase connects to the PostgreSQL database for both services
func ConnectDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=501417 dbname=seller_service port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = database
	// Auto migrate models for both Category and Seller
	if err := DB.AutoMigrate(&models.Category{}, &models.Seller{},models.Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully!")
	return DB, nil
}

// SetupDependencies initializes both the Category and Seller services
func SetupDependencies() (categorypb.CategoryServiceServer, sellerpb.SellerServiceServer, productpb.ProductServiceServer, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, nil, nil, err
	}

	// Initialize CategoryService
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize SellerService
	userRepo := repository.NewUserRepository(db)
	sellerService := service.NewUserService(userRepo)
	sellerHandler := handlers.NewSellerServiceHandler(sellerService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductServiceServer(productRepo)

	return categoryService, sellerHandler, productService, nil
}

// StartGRPCServer starts a single gRPC server for both services
func StartGRPCServer(categoryHandler categorypb.CategoryServiceServer, sellerHandler sellerpb.SellerServiceServer, productHandler productpb.ProductServiceServer) {
	// Create a single gRPC server
	grpcServer := grpc.NewServer()

	// Register both services on the same gRPC server
	categorypb.RegisterCategoryServiceServer(grpcServer, categoryHandler)
	sellerpb.RegisterSellerServiceServer(grpcServer, sellerHandler)
	productpb.RegisterProductServiceServer(grpcServer, productHandler)
	// Listen on port 50052
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Seller and Category service is running on port 50052...")

	// Start serving the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
