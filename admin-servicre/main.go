package main

import (
	"log"
	"net"

	"github.com/faris-muhammed/e-commerce/admin-service/config"
	"github.com/faris-muhammed/e-commerce/admin-service/handlers"
	"github.com/faris-muhammed/e-commerce/admin-service/repository"
	"github.com/faris-muhammed/e-commerce/admin-service/service"
	adminpb "github.com/faris-muhammed/e-protofiles/adminlogin"
	"google.golang.org/grpc"
)

func main() {
	// Database Connection
	db, err := config.ConnectDatabase()
	if err != nil {
		return
	}

	userRepo := repository.NewUserRepository(db)
	adminService := service.NewUserService(userRepo)
	adminHandler := handlers.NewAdminServiceHandler(adminService)
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	adminpb.RegisterAdminServiceServer(grpcServer, adminHandler)

	log.Println("gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
