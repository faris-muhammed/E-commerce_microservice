package repository

import (
	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	UpdateCategory(category *models.Category) error
	DeleteCategory(categoryID int32) error
	GetCategories() ([]models.Category, error)
	FindCategoryByID(id uint) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(categoryID int32) error {
	return r.db.Where("id = ?", categoryID).Delete(&models.Category{}).Error
}

func (r *categoryRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
