package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/faris-muhammed/e-commerce/user-service/helper"
	"github.com/faris-muhammed/e-commerce/user-service/middleware"
	"github.com/faris-muhammed/e-commerce/user-service/models"
	"github.com/faris-muhammed/e-commerce/user-service/repository"
	userpb "github.com/faris-muhammed/e-protofiles/userlogin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("YOURKEY"))

type UserServiceServer struct {
	userpb.UnimplementedLoginServiceServer
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceServer {
	return &UserServiceServer{repo: repo}
}

func (s *UserServiceServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	// Find user by email
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return &userpb.LoginResponse{
			Message: "Invalid email or password",
		}, err
	}

	// Generate token (dummy for now, replace with JWT)
	tokenString, err := middleware.GenerateJWT(uint(user.ID), user.Email, "User")
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		Message: "Login successful",
		Token:   tokenString,
	}, nil
}

func (s *UserServiceServer) Signup(ctx context.Context, req *userpb.SignupRequest) (*userpb.SignupResponse, error) {

	// Check if the email already exists
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err // Something went wrong with the query
	}
	if user != nil {
		return nil, errors.New("email already exists")
	}

	// Generate OTP
	otp := helper.GenerateOTP()
	if err := helper.SendOTP(req.Email, otp); err != nil {
		return nil, errors.New("failed to send OTP")
	}

	// Save OTP details
	otpDetails := models.OTP{
		Email:     req.Email,
		OTP:       otp,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(180 * time.Second),
	}

	if err := s.repo.SaveOrUpdateOTP(otpDetails); err != nil {
		return nil, errors.New("failed to save OTP details")
	}

	// Cache user data
	s.repo.CacheUserData(req.Email, req.Name, req.Password, req.Mobile)

	// Respond with success
	return &userpb.SignupResponse{
		Message: "OTP sent successfully",
		Otp:     otp,
	}, nil
}

func (s *UserServiceServer) VerifyOTP(ctx context.Context, req *userpb.VerifyOTPRequest) (*userpb.VerifyOTPResponse, error) {
	// Fetch the stored OTP details from the repository
	storedOTP, err := s.repo.GetOTPByEmail(req.Email)
	if err != nil {
		return nil, errors.New("OTP not found or expired")
	}

	// Log OTP details for debugging
	fmt.Printf("Stored OTP: %s, User OTP: %s, Expires At: %v\n", storedOTP.OTP, req.Otp, storedOTP.ExpiresAt)

	// Validate OTP
	if storedOTP.OTP != req.Otp {
		return nil, errors.New("invalid OTP")
	}

	// Fetch cached user data from the repository
	userData, err := s.repo.GetCachedUserData(req.Email)
	if err != nil {
		return nil, errors.New("user data not found")
	}

	// Hash the password
	hashedPassword, err := helper.HashPassword(userData.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create the user entity
	user := models.User{
		Name:     userData.Name,
		Email:    userData.Email,
		Mobile:   userData.Mobile,
		Password: hashedPassword,
	}

	// Save the user in the database
	if err := s.repo.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Remove the OTP record
	if err := s.repo.DeleteOTPByEmail(req.Email); err != nil {
		return nil, errors.New("failed to delete OTP record")
	}

	// Respond with success
	return &userpb.VerifyOTPResponse{
		Message: "OTP verified and user created successfully",
	}, nil
}
