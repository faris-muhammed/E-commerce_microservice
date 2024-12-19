package repository

import (
	"errors"

	"github.com/faris-muhammed/e-commerce/admin-service/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(username string) (*models.Admin, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUserByEmail(username string) (*models.Admin, error) {
	var user models.Admin
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
