package repository

import (
	"errors"
	"fmt"

	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) (*models.Product, error)
	GetProductByIDRecover(id uint64) (*models.Product, error)
	GetProductByID(id uint64) (*models.Product, error)
	UpdateProduct(product *models.Product) (*models.Product, error)
	// DeleteProduct(id uint64) error
	ListProducts(sellerid uint64) ([]*models.Product, error)
	GetCategoryByID(categoryID uint64) (*models.Category, error)
	ProductExists(name string, categoryID, sellerID uint64) (bool, error)
	GetAllProducts() ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) (*models.Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) GetProductByIDRecover(id uint64) (*models.Product, error) {
	var product models.Product
	// Log the ID to ensure it's correct
	fmt.Printf("Fetching product with ID: %d\n", id)

	// Check if the product exists and is not soft-deleted in the database
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		// Log the error
		fmt.Printf("Error fetching product: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch product with ID %d: %w", id, err)
	}

	// Return the product if found
	return &product, nil
}
func (r *productRepository) GetProductByID(id uint64) (*models.Product, error) {
	var product models.Product
	// Log the ID to ensure it's correct
	fmt.Printf("Fetching product with ID: %d\n", id)

	// Check if the product exists and is not soft-deleted in the database
	err := r.db.Where("id = ? AND is_deleted=?", id,false).First(&product).Error
	if err != nil {
		// Log the error
		fmt.Printf("Error fetching product: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch product with ID %d: %w", id, err)
	}

	// Return the product if found
	return &product, nil
}
func (r *productRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	
	err := r.db.Where("is_deleted = ?", false).Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	fmt.Printf("Products fetched from DB: %v\n", products)
	return products, nil
}

func (r *productRepository) UpdateProduct(product *models.Product) (*models.Product, error) {
	if err := r.db.Save(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) ListProducts(sellerid uint64) ([]*models.Product, error) {
	var products []*models.Product
	if err := r.db.Where("seller_id=? AND is_deleted = ?", sellerid, false).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetCategoryByID(categoryID uint64) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, categoryID).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
func (r *productRepository) ProductExists(name string, categoryID, sellerID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).
		Where("name = ? AND category_id = ? AND seller_id = ?", name, categoryID, sellerID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
