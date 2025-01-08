package repository

import (
	"errors"
	"fmt"

	"github.com/faris-muhammed/e-commerce/user-service/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	AddToCart(item models.Cart) error
	EditCartItem(userID, productID uint64, quantity uint32) error
	RemoveCartItem(userID uint64) error
	RemoveCartByID(userID, productID uint64) error
	ListCartItems(userID uint64) ([]models.Cart, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddToCart(item models.Cart) error {
	// Check if the item already exists in the cart
	var existing models.Cart
	if err := r.db.Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(&item).Error // Item doesn't exist, create a new one
		}
		return err
	}

	// Update the quantity if the item already exists
	existing.Quantity += item.Quantity
	return r.db.Save(&existing).Error
}

func (r *cartRepository) EditCartItem(userID, productID uint64, quantity uint32) error {
	var item models.Cart
	if err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error; err != nil {
		return err
	}
	item.Quantity = quantity
	return r.db.Save(&item).Error
}

func (r *cartRepository) RemoveCartItem(userID uint64) error {
	// Update the is_deleted field to true instead of deleting the record
	return r.db.Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}

func (r *cartRepository) RemoveCartByID(userID, productID uint64) error {
	// Delete the cart item for the given user and product
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Cart{}).Error
	if err != nil {
		return fmt.Errorf("failed to remove cart item for user %d and product %d: %w", userID, productID, err)
	}
	return nil
}

func (r *cartRepository) ListCartItems(userID uint64) ([]models.Cart, error) {
	var items []models.Cart
	if err := r.db.Where("user_id = ? ", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
