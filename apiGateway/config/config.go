package config

import (
	"github.com/faris-muhammed/e-commerce/apigateway/handlers"
	addresspb "github.com/faris-muhammed/e-protofiles/address"
	adminpb "github.com/faris-muhammed/e-protofiles/adminlogin"
	cartpb "github.com/faris-muhammed/e-protofiles/cart"
	categorypb "github.com/faris-muhammed/e-protofiles/category"
	productpb "github.com/faris-muhammed/e-protofiles/product"
	sellerpb "github.com/faris-muhammed/e-protofiles/sellerlogin"
	userpb "github.com/faris-muhammed/e-protofiles/userlogin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Dependencies struct {
	Router             *gin.Engine
	CategoryHandler    handlers.CategoryHandler
	SellerHandler      handlers.SellerLoginHandler
	ProductHandler     handlers.ProductHandler
	LoginHandler       handlers.AdminLoginHandler
	UserHandler        handlers.UserLoginHandler
	CartHandler        handlers.CartHandler
	UserAddressHandler handlers.UserAddressHandler
}

func InitDependencies() (*Dependencies, func(), error) {
	// Cleanup function to close gRPC connections
	cleanup := func() {}

	// Connect to seller-service
	sellerConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	cleanup = combineCleanup(cleanup, func() { sellerConn.Close() })

	// Initialize handlers for seller-service
	categoryHandler := handlers.NewCategoryHandler(categorypb.NewCategoryServiceClient(sellerConn))
	sellerHandler := handlers.NewSellerLoginHandler(sellerpb.NewSellerServiceClient(sellerConn))
	productHandler := handlers.NewProductHandler(productpb.NewProductServiceClient(sellerConn))

	// Connect to admin-service
	adminConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	cleanup = combineCleanup(cleanup, func() { adminConn.Close() })

	// Initialize handlers for admin-service
	loginHandler := handlers.NewLoginHandler(adminpb.NewAdminServiceClient(adminConn))

	// Connect to user-service
	userConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	cleanup = combineCleanup(cleanup, func() { userConn.Close() })

	// Initialize handlers for user-service
	userHandler := handlers.NewUserLoginHandler(userpb.NewLoginServiceClient(userConn))
	cartHandler := handlers.NewCartHandler(cartpb.NewCartServiceClient(userConn))
	addressHandler := handlers.NewAddressHandler(addresspb.NewUserServiceClient(userConn))

	// Initialize Gin router
	router := gin.Default()
	store := cookie.NewStore([]byte("YOURKEY"))
	router.Use(sessions.Sessions("mysession", store))

	return &Dependencies{
		Router:             router,
		CategoryHandler:    *categoryHandler,
		SellerHandler:      *sellerHandler,
		ProductHandler:     *productHandler,
		LoginHandler:       *loginHandler,
		UserHandler:        *userHandler,
		CartHandler:        *cartHandler,
		UserAddressHandler: *addressHandler,
	}, cleanup, nil
}

// Helper function to combine cleanup functions
func combineCleanup(original, newCleanup func()) func() {
	return func() {
		original()
		newCleanup()
	}
}
