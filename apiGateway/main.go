package main

import (
	"log"

	"github.com/faris-muhammed/e-commerce/apigateway/config"
)

func main() {
	// Initialize dependencies
	deps, cleanup, err := config.InitDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	// Register routes
	config.RegisterRoutes(deps)

	// Start the HTTP server
	log.Println("HTTP Server is running on port 8080...")
	if err := deps.Router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}
