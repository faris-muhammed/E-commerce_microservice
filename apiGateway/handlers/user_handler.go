package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/faris-muhammed/e-commerce/apigateway/middleware"
	userpb "github.com/faris-muhammed/e-protofiles/userlogin"
	"github.com/gin-gonic/gin"
)

type UserLoginHandler struct {
	userService userpb.LoginServiceClient
}

func NewUserLoginHandler(client userpb.LoginServiceClient) *UserLoginHandler {
	return &UserLoginHandler{userService: client}
}

// =============== Signup =======================
func (h *UserLoginHandler) SignUp(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Mobile   string `json:"mobile"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	grpcReq := &userpb.SignupRequest{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := h.userService.Signup(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Generate a session token
	sessionToken, err := middleware.GenerateSessionToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	// Store the session token in an HTTP-only cookie
	c.SetCookie("session_token", sessionToken, 3600, "/", "", false, true)

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{
		"message": resp.Message,
		"otp":     resp.Otp,
	})
}

func (h *UserLoginHandler) VerifyOTP(c *gin.Context) {
	// Retrieve the session token from the cookie
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session token missing or invalid"})
		return
	}

	// Validate the session token and extract the user ID
	userEmail, err := middleware.GetUserIDBySessionToken(sessionToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		OTP string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	grpcReq := &userpb.VerifyOTPRequest{
		Otp:   req.OTP,
		Email: userEmail,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := h.userService.VerifyOTP(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

// UserLogin handles user login requests
func (h *UserLoginHandler) UserLogin(c *gin.Context) {
	// Bind JSON request to struct
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create gRPC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Call the gRPC method
	grpcReq := &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.userService.Login(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate", "err": err.Error()})
		return
	}
	c.SetCookie("jwtTokenUser", resp.Token, int((time.Hour * 6).Seconds()), "/", "", false, true)
	// Return the response to the client
	c.JSON(http.StatusOK, gin.H{
		"message": "User login success!",
		"token":   resp.Token,
	})
}

func (h *UserLoginHandler) UserLogout(c *gin.Context) {
	c.SetCookie("jwtTokenUser", "", -1, "/", "", false, true)
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "User logged out successfully",
		"code":    200,
	})
}
