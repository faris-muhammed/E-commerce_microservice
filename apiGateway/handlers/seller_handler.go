package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	sellerpb "github.com/faris-muhammed/e-protofiles/sellerlogin"
	"github.com/gin-gonic/gin"
)

type SellerLoginHandler struct {
	sellerService sellerpb.SellerServiceClient
}

// NewSellerLoginHandler initializes a new SellerLoginHandler with the gRPC client.
func NewSellerLoginHandler(client sellerpb.SellerServiceClient) *SellerLoginHandler {
	return &SellerLoginHandler{sellerService: client}
}

// LoginHTTP handles seller login via HTTP.
func (h *SellerLoginHandler) LoginHTTP(c *gin.Context) {
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

	grpcReq := &sellerpb.LoginRequest{
		Username: req.Email,
		Password: req.Password,
	}

	// Call gRPC service
	log.Printf("Sending gRPC request: Username=%s, Password=%s", req.Email, req.Password)
	resp, err := h.sellerService.Login(ctx, grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate", "err": err.Error()})
		return
	}
	c.SetCookie("jwtTokenSeller", resp.Token, int((time.Hour * 6).Seconds()), "/", "", false, true)
	// Return the response with the token
	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
		"token":   resp.Token,
	})
}


// LogoutHTTP handles seller logout via HTTP.
func (h *SellerLoginHandler) LogoutHTTP(c *gin.Context) {
	c.SetCookie("jwtTokenSeller", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Seller logged out successfully",
	})
}
