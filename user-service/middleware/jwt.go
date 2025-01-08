package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("SECRETKEY"))

func GenerateJWT(userId uint, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"userId":   userId,
		"role":     role,                                  // Include the user's role
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
