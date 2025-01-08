package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	sellerpb "github.com/faris-muhammed/e-protofiles/category"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	sellerService sellerpb.CategoryServiceClient
}

func NewCategoryHandler(sellerClient sellerpb.CategoryServiceClient) *CategoryHandler {
	return &CategoryHandler{
		sellerService: sellerClient,
	}
}

func (h *CategoryHandler) AddCategoryHTTP(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward to seller-service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sellerReq := &sellerpb.AddCategoryRequest{
		Name: req.Name,
		Description:  req.Description,
	}
	_, err := h.sellerService.AddCategory(ctx, sellerReq)
	if err != nil {
		fmt.Printf("Failed to add category %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add category in seller-service."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category added successfully in seller-service!"})
}
func (h *CategoryHandler) EditCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sellerReq := &sellerpb.EditCategoryRequest{
		CategoryId:   uint64(id),
		Name: req.Name,
		Description:  req.Description,
	}
	_, err = h.sellerService.EditCategory(ctx, sellerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category in seller-service."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully!"})
}

func (h *CategoryHandler) GetCategoriesHTTP(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sellerReq := &sellerpb.Empty{}
	resp, err := h.sellerService.GetCategories(ctx, sellerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories from seller-service."})
		return
	}
	c.JSON(http.StatusOK, resp.Categories)
}
func (h *CategoryHandler) DeleteCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Forward to seller-service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sellerReq := &sellerpb.DeleteCategoryRequest{
		CategoryId: uint64(id),
	}
	resp, err := h.sellerService.DeleteCategory(ctx, sellerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category in seller-service."})
		return
	}
	// Check response and return appropriate HTTP message
	if resp.Message == "Category not found!" {
		c.JSON(http.StatusNotFound, gin.H{"message": "Category not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully in seller-service!"})
}
