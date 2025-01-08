package routes

import (
	"github.com/faris-muhammed/e-commerce/apigateway/handlers"
	"github.com/faris-muhammed/e-commerce/apigateway/middleware"
	"github.com/gin-gonic/gin"
)

func CategoryGroup(h *handlers.CategoryHandler, r *gin.RouterGroup) {
	r.POST("/categories", middleware.AuthMiddleware("Admin"), h.AddCategoryHTTP)
	r.GET("/categories", middleware.AuthMiddleware("Admin"), h.GetCategoriesHTTP)
	r.PUT("/categories/:id", middleware.AuthMiddleware("Admin"), h.EditCategoryHTTP)
	r.DELETE("/categories/:id/", middleware.AuthMiddleware("Admin"), h.DeleteCategoryHTTP)
}
func LoginGroup(h *handlers.AdminLoginHandler, r *gin.RouterGroup) {
	r.GET("/login", h.LoginHTTP)
	r.DELETE("/logout", middleware.AuthMiddleware("Admin"), h.LogoutHTTP)
}
func SellerGroup(h *handlers.SellerLoginHandler, r *gin.RouterGroup) {
	r.GET("/login", h.LoginHTTP)
	r.DELETE("/logout", middleware.AuthMiddleware("Seller"), h.LogoutHTTP)
}
func UserGroup(h *handlers.UserLoginHandler, r *gin.RouterGroup) {
	r.POST("/signup", h.SignUp)
	r.POST("/verify-otp", h.VerifyOTP)
	r.GET("/login", h.UserLogin)
	r.DELETE("/logout", middleware.AuthMiddleware("User"), h.UserLogout)
}
func ProductGroup(h *handlers.ProductHandler, r *gin.RouterGroup) {
	r.POST("/product", middleware.AuthMiddleware("Seller"), h.CreateProduct)
	r.PUT("/product/:id", middleware.AuthMiddleware("Seller"), h.UpdateProduct)
	r.GET("/productslr", middleware.AuthMiddleware("Seller"), h.ListProductsSeller)
	r.GET("/product/:id", middleware.AuthMiddleware("Seller"), h.GetProduct)
	r.PUT("/product/delete/:id", middleware.AuthMiddleware("Seller"), h.DeleteProduct)
	r.GET("/product", middleware.AuthMiddleware("Seller"), h.ListProducts)
}

func CartGroup(h *handlers.CartHandler, r *gin.RouterGroup) {
	r.GET("/cart", middleware.AuthMiddleware("User"), h.ListCartItems)
	r.POST("/cart", middleware.AuthMiddleware("User"), h.AddToCart)
	r.PUT("/cart", middleware.AuthMiddleware("User"), h.EditCartItem)
	r.DELETE("/cart/delete", middleware.AuthMiddleware("User"), h.RemoveCart)
	r.PUT("/cart/:id", middleware.AuthMiddleware("User"), h.RemoveCartByID)
}

func AddressGroup(h *handlers.UserAddressHandler, r *gin.RouterGroup) {
	r.GET("/address", middleware.AuthMiddleware("User"), h.GetUserAddresses)
	r.POST("/address", middleware.AuthMiddleware("User"), h.AddUserAddress)
	r.PUT("/address/:id", middleware.AuthMiddleware("User"), h.EditUserAddress)
	r.DELETE("/address", middleware.AuthMiddleware("User"), h.DeleteUserAddress)
}
