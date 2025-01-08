package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/faris-muhammed/e-commerce/user-service/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(email string) (*models.User, error)
	SaveOrUpdateOTP(details models.OTP) error
	GetOTPByEmail(email string) (*models.OTP, error)
	DeleteOTPByEmail(email string) error
	CacheUserData(email, name, password, mobile string)
	GetCachedUserData(email string) (*models.CachedUser, error)
	CreateUser(user models.User) error
}

type userRepository struct {
	db    *gorm.DB
	cache map[string]models.CachedUser
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db, cache: make(map[string]models.CachedUser)}
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("The error email", err, email)
			return nil, nil
		}
		fmt.Println("Error while querying user:", err)
		return nil, err
	}
	fmt.Println("User found with email:", email)
	return &user, nil
}
func (r *userRepository) SaveOrUpdateOTP(details models.OTP) error {
    var existing models.OTP
    // First, check if any OTP exists for this email
    if err := r.db.Where("email = ?", details.Email).First(&existing).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return r.db.Create(&details).Error // No existing OTP found, create a new one
        }
        return err
    }

    // If the OTP exists and it's expired or invalid, delete it
    if existing.ExpiresAt.Before(time.Now()) || existing.OTP != details.OTP {
        if err := r.db.Delete(&existing).Error; err != nil {
            return err // Error while deleting old OTP
        }
    }

    // Insert the new OTP as the latest one
    return r.db.Create(&details).Error
}



func (r *userRepository) GetOTPByEmail(email string) (*models.OTP, error) {
	var otp models.OTP
	if err := r.db.Where("email = ?", email).First(&otp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("otp not found")
		}
		return nil, err
	}
	
	// Check if the OTP is expired
	if otp.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("otp has expired")
	}

	return &otp, nil
}



func (r *userRepository) DeleteOTPByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.OTP{}).Error
}

func (r *userRepository) CacheUserData(email, name, password, mobile string) {
	r.cache[email] = models.CachedUser{Name: name, Email: email, Password: password, Mobile: mobile}
}

func (r *userRepository) GetCachedUserData(email string) (*models.CachedUser, error) {
	user, exists := r.cache[email]
	if !exists {
		return nil, errors.New("user data not found in cache")
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user models.User) error {
	return r.db.Create(&user).Error
}
