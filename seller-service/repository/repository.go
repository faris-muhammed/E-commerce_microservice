package repository

import (
	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories() ([]*models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	DeleteCategory(id uint) error
	UpdateCategory(category *models.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetCategories() ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	if err := r.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) DeleteCategory(id uint) error {
	if err := r.db.Delete(&models.Category{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) UpdateCategory(category *models.Category) error {
	if err := r.db.Save(&category).Error; err != nil {
		return err
	}
	return nil
}
