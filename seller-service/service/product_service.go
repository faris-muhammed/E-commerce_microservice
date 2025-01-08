package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"github.com/faris-muhammed/e-commerce/seller-service/repository"
	productpb "github.com/faris-muhammed/e-protofiles/product"
	"gorm.io/gorm"
)

type ProductServiceServer struct {
	productpb.UnimplementedProductServiceServer
	repo repository.ProductRepository
}

func NewProductServiceServer(repo repository.ProductRepository) *ProductServiceServer {
	return &ProductServiceServer{repo: repo}
}
func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	// Check if the category exists
	category, err := s.repo.GetCategoryByID(req.CategoryId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category with ID %d does not exist", req.CategoryId)
		}
		return nil, fmt.Errorf("failed to check category existence: %w", err)
	}

	// Check if the product already exists
	exists, err := s.repo.ProductExists(req.Name, req.CategoryId, req.SellerId)
	if err != nil {
		return nil, fmt.Errorf("failed to check product existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("product with name '%s' already exists in category ID %d for seller ID %d", req.Name, req.CategoryId, req.SellerId)
	}

	fmt.Println("Stock in seller-service",req.Stock)
	// Proceed to create the product if it does not exist
	product := &models.Product{
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: uint64(category.ID),
		SellerID:   req.SellerId,
	}

	createdProduct, err := s.repo.CreateProduct(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &productpb.CreateProductResponse{
		Message: "Product created successfully",
		Product: &productpb.Product{
			Id:         uint64(createdProduct.ID),
			Name:       createdProduct.Name,
			Price:      createdProduct.Price,
			Stock:      createdProduct.Stock,
			CategoryId: createdProduct.CategoryID,
			SellerId:   uint64(createdProduct.SellerID),
			IsDeleted:  createdProduct.IsDeleted,
		},
	}, nil
}

func (s *ProductServiceServer) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	product, err := s.repo.GetProductByID(uint64(req.ProductId))
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	fmt.Println("Product Stock", product.Stock)
	return &productpb.GetProductResponse{
		Product: &productpb.Product{
			Id:         uint64(product.ID),
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryId: product.CategoryID,
		},
	}, nil
}
func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
	// Retrieve the existing product from the repository
	existingProduct, err := s.repo.GetProductByID(req.Id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Update only the fields that are provided in the request
	if req.Name != "" {
		existingProduct.Name = req.Name
	}
	if req.Price != 0 {
		existingProduct.Price = req.Price
	}
	if req.Stock != 0 {
		existingProduct.Stock = req.Stock
	}
	if req.CategoryId != 0 {
		existingProduct.CategoryID = req.CategoryId
	}

	// Save the updated product to the repository
	updatedProduct, err := s.repo.UpdateProduct(existingProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &productpb.UpdateProductResponse{
		Message: "Product updated successfully",
		Product: &productpb.Product{
			Id:         uint64(updatedProduct.ID),
			Name:       updatedProduct.Name,
			Price:      updatedProduct.Price,
			Stock:      updatedProduct.Stock,
			CategoryId: updatedProduct.CategoryID,
		},
	}, nil
}
func (s *ProductServiceServer) ModifyProductStatus(ctx context.Context, req *productpb.ModifyProductStatusRequest) (*productpb.ModifyProductStatusResponse, error) {
	// Check if the product exists
	product, err := s.repo.GetProductByIDRecover(req.ProductId)
	if err != nil {
		return nil, err
	}

	// Check if action is provided
	if req.Action == "" {
		return nil, errors.New("action is required")
	}

	// Based on the action, modify the product status
	switch req.Action {
	case "soft_delete":
		// Soft delete: set is_deleted to true
		product.IsDeleted = true
		_, err := s.repo.UpdateProduct(product) // Handle both return values
		if err != nil {
			return nil, err
		}
		return &productpb.ModifyProductStatusResponse{Message: "Product soft deleted successfully"}, nil

	case "recover":
		// Recover: set is_deleted to false
		product.IsDeleted = false
		_, err := s.repo.UpdateProduct(product)
		if err != nil {
			return nil, err
		}
		return &productpb.ModifyProductStatusResponse{Message: "Product recovered successfully"}, nil

	default:
		return nil, errors.New("invalid action, must be 'soft_delete' or 'recover'")
	}
}

func (s *ProductServiceServer) ListProductsSeller(ctx context.Context, req *productpb.ListProductsSellerRequest) (*productpb.ListProductsSellerResponse, error) {
	products, err := s.repo.ListProducts(req.SellerId)
	if err != nil {
		return nil, err
	}

	var grpcProducts []*productpb.Product
	for _, product := range products {
		fmt.Printf("Product ID: %d, Name: %s\n", product.ID, product.Name)
		grpcProducts = append(grpcProducts, &productpb.Product{
			Id:         uint64(product.ID),
			Name:       product.Name,
			Price:      product.Price,
			CategoryId: product.CategoryID,
			SellerId:   uint64(product.SellerID),
			IsDeleted:  product.IsDeleted,
		})
	}

	return &productpb.ListProductsSellerResponse{Products: grpcProducts}, nil
}
func (s *ProductServiceServer) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	var grpcProducts []*productpb.Product
	for _, product := range products {
		fmt.Printf("Product ID: %d, Name: %s\n", product.ID, product.Name)
		grpcProducts = append(grpcProducts, &productpb.Product{
			Id:         uint64(product.ID),
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryId: product.CategoryID,
		})
	}
	fmt.Printf("Seller-Service Response: %+v\n", grpcProducts)
	return &productpb.ListProductsResponse{Products: grpcProducts}, nil
}
