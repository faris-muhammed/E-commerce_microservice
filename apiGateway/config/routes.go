package config

import (
	"github.com/faris-muhammed/e-commerce/apigateway/routes"
)

func RegisterRoutes(deps *Dependencies) {
	// Admin routes
	categoryGroup := deps.Router.Group("/admin")
	routes.CategoryGroup(&deps.CategoryHandler, categoryGroup)

	// Seller routes
	sellerGroup := deps.Router.Group("/seller")
	routes.SellerGroup(&deps.SellerHandler, sellerGroup)

	productGroup := deps.Router.Group("/seller")
	routes.ProductGroup(&deps.ProductHandler, productGroup)

	// Login routes
	loginGroup := deps.Router.Group("/")
	routes.LoginGroup(&deps.LoginHandler, loginGroup)

	// User routes
	userGroup := deps.Router.Group("/user")
	routes.UserGroup(&deps.UserHandler, userGroup)

	// Cart routes
	cartGroup := deps.Router.Group("/user")
	routes.CartGroup(&deps.CartHandler, cartGroup)

	// Address routes
	addressGroup := deps.Router.Group("/user")
	routes.AddressGroup(&deps.UserAddressHandler, addressGroup)
}
