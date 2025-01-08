package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/faris-muhammed/e-commerce/seller-service/repository"
	"github.com/golang-jwt/jwt/v4"
)

type UserService interface {
	Login(email, password string) (string, error)
}

type userService struct {
	repo repository.SellerRepository
}

func NewUserService(repo repository.SellerRepository) UserService {
	return &userService{repo: repo}
}

var Role ="Seller"
var jwtKey = []byte(os.Getenv("SECRETKEY")) // Ensure you store this securely, ideally in an environment variable

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Here, you'd compare the password hash with bcrypt or other secure hashing methods
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	 return "", errors.New("invalid email or password")
	// }

	// Generate JWT token with the user's role
	tokenString, err := generateJWT(uint(user.ID), user.Username, Role)
	if err != nil {
		return "", err
	}
	fmt.Println("Generated JWT:", tokenString)
	return tokenString, nil
}

// Generate JWT token with role
func generateJWT(userId uint, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"userId":   userId,
		"role":     role, // Include the user's role
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
