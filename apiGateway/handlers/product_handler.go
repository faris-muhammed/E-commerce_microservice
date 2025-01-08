package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	productpb "github.com/faris-muhammed/e-protofiles/product"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService productpb.ProductServiceClient
}

func NewProductHandler(client productpb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{
		productService: client,
	}
}

// CreateProduct handles the creation of a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	sellerID := c.GetUint("userid")
	var req struct {
		Name       string  `json:"name"`
		Price      float64 `json:"price"`
		Stock      uint64  `json:"stock"`
		CategoryId uint64  `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Stock in api gateway",req.Stock)
	grpcReq := &productpb.CreateProductRequest{
		SellerId:   uint64(sellerID),
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryId: req.CategoryId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.productService.CreateProduct(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProduct handles fetching a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	req := &productpb.GetProductRequest{
		ProductId: uint64(id),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.productService.GetProduct(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

type updateProductRequest struct {
	ID         uint64   `json:"id"`
	Name       *string  `json:"name,omitempty"`
	Price      *float64 `json:"price,omitempty"`
	Stock      *uint64  `json:"stock,omitempty"`
	CategoryID *uint64  `json:"category_id,omitempty"`
}

// UpdateProduct handles updating an existing product
// func (h *ProductHandler) UpdateProduct(c *gin.Context) {
//     idstr := c.Param("id")
//     productID, err := strconv.Atoi(idstr)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
//         return
//     }

//     // Bind the request body to the updateProductRequest
//     var input updateProductRequest
//     if err := c.ShouldBindJSON(&input); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
//         return
//     }

//     // Set the product ID in the request
//     input.ID = uint64(productID)

//     // Convert the input to the gRPC request object
//     grpcReq := &productpb.UpdateProductRequest{
//         Id:         input.ID,
//         Name:       *input.Name,
//         Price:      *input.Price,
//         Stock:      *input.Stock,
//         CategoryId: *input.CategoryID,
//     }

//     // Forward the request to the gRPC client to update the product
//     updatedProduct, err := h.productService.UpdateProduct(c, grpcReq)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     // Return the updated product as a response
//     c.JSON(http.StatusOK, updatedProduct)
// }

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idstr := c.Param("id")
	productID, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Bind the request body to the updateProductRequest
	var input updateProductRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Set the product ID in the request
	input.ID = uint64(productID)

	// Convert the input to the gRPC request object
	grpcReq := &productpb.UpdateProductRequest{
		Id: input.ID,
	}

	if input.Name != nil {
		grpcReq.Name = *input.Name
	}
	if input.Price != nil {
		grpcReq.Price = *input.Price
	}
	if input.Stock != nil {
		grpcReq.Stock = *input.Stock
	}
	if input.CategoryID != nil {
		grpcReq.CategoryId = *input.CategoryID
	}

	// Forward the request to the gRPC client to update the product
	updatedProduct, err := h.productService.UpdateProduct(c, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated product as a response
	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct handles soft-deleting a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Get product ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Get the action (soft_delete or recover) from query parameter
	action := c.DefaultQuery("action", "")
	if action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action is required (must be 'soft_delete' or 'recover')"})
		return
	}

	// Prepare the request with product ID and action
	req := &productpb.ModifyProductStatusRequest{
		ProductId: uint64(id),
		Action:    action,
	}

	// Call the ModifyProductStatus service method
	resp, err := h.productService.ModifyProductStatus(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response from service
	c.JSON(http.StatusOK, resp)
}

// ListProducts handles listing all products, optionally including soft-deleted ones
func (h *ProductHandler) ListProductsSeller(c *gin.Context) {
	id := c.GetUint("userid")

	req := &productpb.ListProductsSellerRequest{
		SellerId: uint64(id),
	}

	resp, err := h.productService.ListProductsSeller(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Response: %+v\n", resp)
	c.JSON(http.StatusOK, resp)
}
