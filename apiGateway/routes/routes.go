package routes

import (
	"github.com/faris-muhammed/e-commerce/apigateway/handlers"
	"github.com/faris-muhammed/e-commerce/apigateway/middleware"
	"github.com/gin-gonic/gin"
)

// var Role = "Admin"

func CategoryGroup(h *handlers.CategoryHandler, r *gin.RouterGroup) {
	r.POST("/categories", middleware.AuthMiddleware("Admin"), h.AddCategoryHTTP)
	r.GET("/categories", middleware.AuthMiddleware("Admin"), h.GetCategoriesHTTP)
	r.PUT("/categories/:id", middleware.AuthMiddleware("Admin"), h.EditCategoryHTTP)
	r.DELETE("/categories/:id/", middleware.AuthMiddleware("Admin"), h.DeleteCategoryHTTP)
}
func LoginGroup(h *handlers.LoginHandler, r *gin.RouterGroup) {
	r.GET("/login", h.LoginHTTP)
	r.DELETE("/logout", middleware.AuthMiddleware("Admin"), h.LogoutHTTP)
}
