package service

import (
	"context"
	"errors"

	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"github.com/faris-muhammed/e-commerce/seller-service/repository"
	seller "github.com/faris-muhammed/e-protofiles/category"
)

type categoryService struct {
	seller.UnimplementedCategoryServiceServer // Embed the unimplemented server
	repo repository.CategoryRepository
}

// NewCategoryService initializes a new CategoryService.
func NewCategoryService(repo repository.CategoryRepository) seller.CategoryServiceServer {
	return &categoryService{repo: repo}
}

func (s *categoryService) AddCategory(ctx context.Context, req *seller.AddCategoryRequest) (*seller.CategoryResponse, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	_, err := s.repo.CreateCategory(category)
	if err != nil {
		return nil, err
	}
	return &seller.CategoryResponse{Message: "Category added successfully!"}, nil
}

func (s *categoryService) EditCategory(ctx context.Context, req *seller.EditCategoryRequest) (*seller.CategoryResponse, error) {
	category, err := s.repo.GetCategoryByID(uint(req.CategoryId))
	if err != nil {
		return nil, errors.New("category not found")
	}

	category.Name = req.Name
	category.Description = req.Description

	if err := s.repo.UpdateCategory(category); err != nil {
		return nil, err
	}

	return &seller.CategoryResponse{Message: "Category updated successfully!"}, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, req *seller.DeleteCategoryRequest) (*seller.CategoryResponse, error) {
	if err := s.repo.DeleteCategory(uint(req.CategoryId)); err != nil {
		return nil, errors.New("category not found or deletion failed")
	}

	return &seller.CategoryResponse{Message: "Category deleted successfully!"}, nil
}

func (s *categoryService) GetCategories(ctx context.Context, req *seller.Empty) (*seller.CategoryListResponse, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, err
	}

	var categoryList []*seller.Category
	for _, c := range categories {
		categoryList = append(categoryList, &seller.Category{
			Id:          uint64(c.ID),
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return &seller.CategoryListResponse{Categories: categoryList}, nil
}
