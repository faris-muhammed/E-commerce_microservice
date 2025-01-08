package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cartpb "github.com/faris-muhammed/e-protofiles/cart"
	productpb "github.com/faris-muhammed/e-protofiles/product"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService cartpb.CartServiceClient
}

func NewCartHandler(client cartpb.CartServiceClient) *CartHandler {
	return &CartHandler{cartService: client}
}

// AddToCart adds a product to the user's cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("userid")
	var req struct {
		ProductID uint64 `json:"product_id" binding:"required"`
		Quantity  uint32 `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &cartpb.AddToCartRequest{
		UserId:    uint64(userID),
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.cartService.AddToCart(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": resp.Message,
	})
}

// ListCartItems lists all items in the user's cart
func (h *CartHandler) ListCartItems(c *gin.Context) {
	userID := c.GetUint("userid")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcReq := &cartpb.ListCartRequest{UserId: uint64(userID)}
	resp, err := h.cartService.ListCart(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cart_items": resp.Items,
		"total":      resp.TotalAmount,
	})
}

// RemoveCart removes all items from the user's cart
func (h *CartHandler) RemoveCart(c *gin.Context) {
	userID := c.GetUint("userid")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcReq := &cartpb.RemoveCartRequest{UserId: uint64(userID)}
	resp, err := h.cartService.RemoveCart(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *CartHandler) RemoveCartByID(c *gin.Context) {
	// Parse the user_id and product_id from the request parameters
	userID := c.GetUint("userid")         // User ID in the URL path
	productIDStr := c.Param("id") // Product ID in the URL path

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Create the gRPC request with the parsed user_id and product_id
	req := &cartpb.RemoveCartByIDRequest{
		UserId:    uint64(userID),
		ProductId: productID,
	}

	// Call the service to remove the cart item
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.cartService.RemoveCartByID(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the success message
	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func (h *CartHandler) EditCartItem(c *gin.Context) {
	userID := c.GetUint("userid")
	var req struct {
		ProductID uint64 `json:"product_id" binding:"required"`
		Quantity  uint32 `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &cartpb.EditCartRequest{
		UserId:    uint64(userID),
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.cartService.EditCart(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

// =============== List Products ====================
func (h *ProductHandler) ListProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := h.productService.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		fmt.Printf("Error from cart service: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": resp.Products,
	})
}
