package repository

import (
	"errors"

	"github.com/faris-muhammed/e-commerce/seller-service/models"
	"gorm.io/gorm"
)

type SellerRepository interface {
	FindUserByEmail(username string) (*models.Seller, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) SellerRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUserByEmail(username string) (*models.Seller, error) {
	var user models.Seller
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
