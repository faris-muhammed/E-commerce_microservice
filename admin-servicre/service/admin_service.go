package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/faris-muhammed/e-commerce/admin-service/repository"
	adminpb "github.com/faris-muhammed/e-protofiles/adminlogin"
	"github.com/golang-jwt/jwt/v4"
)

type AdminServiceServer struct {
	adminpb.UnimplementedAdminServiceServer
	repo repository.AdminRepository
}


func NewUserService(repo repository.AdminRepository) *AdminServiceServer {
	return &AdminServiceServer{repo: repo}
}

var Role = "Admin"

var jwtKey = []byte(os.Getenv("SECRETKEY"))

func (s *AdminServiceServer) Login(ctx context.Context, req *adminpb.LoginRequest) (*adminpb.LoginResponse, error) {	
    user, err := s.repo.FindUserByEmail(req.Username)
    if err != nil {
        return nil, errors.New("invalid email or password")
    }

    // Compare password securely (use bcrypt or another secure method here)
    // if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
    //     return nil, errors.New("invalid email or password")
    // }

    // Generate JWT token
    tokenString, err := generateJWT(uint(user.Id), user.Username, Role)
    if err != nil {
        return nil, err
    }

    fmt.Println("Generated JWT:", tokenString)
    return &adminpb.LoginResponse{Message: "Login Successful", Token: tokenString}, nil
}


// Generate JWT token with role
func generateJWT(userId uint, username, role string) (string, error) {
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
