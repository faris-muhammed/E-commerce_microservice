package service

import (
	"errors"

	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"github.com/faris-muhammed/e-commerce/seller-service/repository"
	seller "github.com/faris-muhammed/e-protofiles/category"
)

type CategoryService interface {
	AddCategory(req *seller.AddCategoryRequest) (*seller.CategoryResponse, error)
	EditCategory(req *seller.EditCategoryRequest) (*seller.CategoryResponse, error)
	DeleteCategory(req *seller.DeleteCategoryRequest) (*seller.CategoryResponse, error)
	GetCategories() (*seller.CategoryListResponse, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService initializes a new CategoryService.
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) AddCategory(req *seller.AddCategoryRequest) (*seller.CategoryResponse, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}
	return &seller.CategoryResponse{Message: "Category added successfully!"}, nil
}

func (s *categoryService) EditCategory(req *seller.EditCategoryRequest) (*seller.CategoryResponse, error) {
	category, err := s.repo.FindCategoryByID(uint(req.CategoryId))
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

//	func (s *categoryService) DeleteCategory(req *seller.DeleteCategoryRequest) (*seller.CategoryResponse, error) {
//		if err := s.repo.DeleteCategory(uint(req.CategoryId)); err != nil {
//			return nil, errors.New("category not found or deletion failed")
//		}
//		return &seller.CategoryResponse{Message: "Category deleted successfully!"}, nil
//	}
func (s *categoryService) DeleteCategory(req *seller.DeleteCategoryRequest) (*seller.CategoryResponse, error) {
	// Check if the categoryId is non-negative
	if req.CategoryId < 0 {
		return nil, errors.New("invalid category ID: must be a non-negative integer")
	}

	// Call repository to delete the category
	if err := s.repo.DeleteCategory(req.CategoryId); err != nil {
		return nil, errors.New("category not found or deletion failed")
	}

	return &seller.CategoryResponse{Message: "Category deleted successfully!"}, nil
}

func (s *categoryService) GetCategories() (*seller.CategoryListResponse, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, err
	}

	var categoryList []*seller.Category
	for _, c := range categories {
		categoryList = append(categoryList, &seller.Category{
			Id:          c.ID,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return &seller.CategoryListResponse{Categories: categoryList}, nil
}
