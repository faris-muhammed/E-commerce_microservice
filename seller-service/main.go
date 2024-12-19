package main

import (
	"log"

	"github.com/faris-muhammed/e-commerce/seller-service/config"
)

func main() {
	// Step 1: Initialize dependencies
	handler, err := config.SetupDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Step 2: Start gRPC server
	config.StartGRPCServer(handler)
}
