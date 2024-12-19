package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	adminpb "github.com/faris-muhammed/e-protofiles/adminlogin"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	adminClient adminpb.AdminServiceClient
}

func NewLoginHandler(client adminpb.AdminServiceClient) *LoginHandler {
	return &LoginHandler{adminClient: client}
}

func (h *LoginHandler) LoginHTTP(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Create gRPC request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &adminpb.LoginRequest{
		Username: req.Email,
		Password: req.Password,
	}

	// Call gRPC service
	log.Printf("Sending gRPC request: Username=%s, Password=%s", req.Email, req.Password)
	resp, err := h.adminClient.Login(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate", "err": err.Error()})
		return
	}
	c.SetCookie("jwtTokenAdmin", resp.Token, int((time.Hour * 6).Seconds()), "/", "", false, true)
	// Return the response with the token
	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
		"token":   resp.Token,
	})
}

func (h *LoginHandler) LogoutHTTP(c *gin.Context) {
	c.SetCookie("jwtTokenAdmin", "", -1, "/", "", false, true)
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Admin logged out successfully",
		"code":    200,
	})
}
