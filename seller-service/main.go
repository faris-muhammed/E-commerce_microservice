package main

import (
	"log"

	"github.com/faris-muhammed/e-commerce/seller-service/config"
)

func main() {
	// Initialize and set up dependencies for both Category and Seller services
	categoryHandler, sellerHandler, productHandler, err := config.SetupDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Start the gRPC server for both services on port 50051
	config.StartGRPCServer(categoryHandler, sellerHandler, productHandler)
}
