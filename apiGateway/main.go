package main

import (
	"log"

	"github.com/faris-muhammed/e-commerce/apigateway/handlers"
	"github.com/faris-muhammed/e-commerce/apigateway/routes"
	adminpb "github.com/faris-muhammed/e-protofiles/adminLogin"
	sellerpb "github.com/faris-muhammed/e-protofiles/seller"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Initialize Gin Router
	r := gin.Default()

	//===========================================================
	// gRPC connection to seller-service
	seller, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to seller-service: %v", err)
	}
	defer seller.Close()

	// gRPC Client for Seller Service
	sellerClient := sellerpb.NewCategoryServiceClient(seller)
	// Pass the client to the CategoryHandler
	categoryHandler := handlers.NewCategoryHandler(sellerClient)

	// Grouping and route setup
	categoryGroup := r.Group("/admin")
	routes.CategoryGroup(categoryHandler, categoryGroup)
	//===========================================================

	//===========================================================
	// gRPC connection to seller-service
	admin, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to seller-service: %v", err)
	}
	defer admin.Close()
	// gRPC Client for admin-service
	adminClient := adminpb.NewAdminServiceClient(admin)
	// Pass the client to the LoginHandler
	loginHandler := handlers.NewLoginHandler(adminClient)
	// Grouping and route setup
	loginGroup := r.Group("/")
	routes.LoginGroup(loginHandler, loginGroup)
	//==========================================================

	log.Println("HTTP Server is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}
