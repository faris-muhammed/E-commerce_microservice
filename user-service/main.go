package main

import (
	"log"
	"net"

	"github.com/faris-muhammed/e-commerce/user-service/config"
	"github.com/faris-muhammed/e-commerce/user-service/repository"
	"github.com/faris-muhammed/e-commerce/user-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	addresspb "github.com/faris-muhammed/e-protofiles/address"
	cartpb "github.com/faris-muhammed/e-protofiles/cart"
	productpb "github.com/faris-muhammed/e-protofiles/product"
	userpb "github.com/faris-muhammed/e-protofiles/userlogin"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {
	// Database Connection
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
		return
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	cartRepo := repository.NewCartRepository(db)
	addressRepo := repository.NewUserAddressRepository(db)

	// Setup gRPC connection to Seller Service (Product Service)
	sellerConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to Seller service: %v", err)
	}
	defer sellerConn.Close()

	// Create a ProductServiceClient
	productService := productpb.NewProductServiceClient(sellerConn)
	// Create CartService with repository and ProductService client
	cartService := service.NewCartService(cartRepo, productService)
	addressService := service.NewUserAddressService(addressRepo)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen on port 50053: %v", err)
	}

	// Setup the gRPC server
	grpcServer := grpc.NewServer()
	userpb.RegisterLoginServiceServer(grpcServer, userService)
	cartpb.RegisterCartServiceServer(grpcServer, cartService)
	addresspb.RegisterUserServiceServer(grpcServer, addressService)

	log.Println("gRPC Server is running on port 50053...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
