package repository

import (
	"fmt"

	"github.com/faris-muhammed/e-commerce/user-service/models"
	"gorm.io/gorm"
)

// UserRepository interface defines methods for interacting with the user data
type UserAddressRepository interface {
	CreateUserAddress(address *models.Address) error
	GetUserAddresses(userID uint64) ([]models.Address, error)
	GetAddressByID(addressID uint64) (*models.Address, error)
	GetAddressByStreet(street string) (*models.Address, error)
	UpdateUserAddress(address *models.Address) (*models.Address, error)
	DeleteUserAddress(address *models.Address) error
}

// userRepo is a concrete implementation of UserRepository
type userRepo struct {
	db *gorm.DB
}

func NewUserAddressRepository(db *gorm.DB) UserAddressRepository {
	return &userRepo{db: db}
}

// CreateUserAddress creates a new user address in the database
func (r *userRepo) CreateUserAddress(address *models.Address) error {
	result := r.db.Create(address)
	if result.Error != nil {
		return fmt.Errorf("failed to create user address: %w", result.Error)
	}
	return nil
}

// GetUserAddresses retrieves all addresses for a given user
func (r *userRepo) GetUserAddresses(userID uint64) ([]models.Address, error) {
	var addresses []models.Address
	result := r.db.Where("user_id = ?", userID).Find(&addresses)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch addresses for user %d: %w", userID, result.Error)
	}
	return addresses, nil
}

// GetAddressByID retrieves a specific address by its ID
func (r *userRepo) GetAddressByID(addressID uint64) (*models.Address, error) {
	var address models.Address
	result := r.db.First(&address, addressID)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch address with ID %d: %w", addressID, result.Error)
	}
	return &address, nil
}

func (r *userRepo) GetAddressByStreet(street string) (*models.Address, error) {
	var address models.Address
	// Use "WHERE" condition to find the record by street
	err := r.db.Where("street = ?", street).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil 
}

// UpdateUserAddress updates an existing user address
func (r *userRepo) UpdateUserAddress(address *models.Address) (*models.Address, error) {
	if err := r.db.Save(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

// DeleteUserAddress deletes a user address from the database
func (r *userRepo) DeleteUserAddress(address *models.Address) error {
	result := r.db.Delete(address)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user address: %w", result.Error)
	}
	return nil
}
