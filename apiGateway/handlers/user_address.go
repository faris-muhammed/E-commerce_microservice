package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	addresspb "github.com/faris-muhammed/e-protofiles/address"
	"github.com/gin-gonic/gin"
)

// UserAddressHandler handles HTTP requests for user addresses
type UserAddressHandler struct {
	userService addresspb.UserServiceClient
}

func NewAddressHandler(client addresspb.UserServiceClient) *UserAddressHandler {
	return &UserAddressHandler{userService: client}
}

// AddUserAddress HTTP handler to add a user address
func (h *UserAddressHandler) AddUserAddress(c *gin.Context) {
	userID := c.GetUint("userid")
	var request struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		State      string `json:"state"`
		PostalCode string `json:"postal_code"`
		Country    string `json:"country"`
	}

	// Bind the request JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	fmt.Println("Request", request, userID)
	// Call the AddUserAddress gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userService.AddUserAddress(ctx, &addresspb.AddUserAddressRequest{
		UserId:     uint64(userID),
		Street:     request.Street,
		City:       request.City,
		State:      request.State,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    resp.Message,
		"address_id": resp.AddressId,
	})
}

// GetUserAddresses HTTP handler to get all user addresses
func (h *UserAddressHandler) GetUserAddresses(c *gin.Context) {
	userID := c.GetUint("userid")

	// Call the GetUserAddresses gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userService.GetUserAddresses(ctx, &addresspb.GetUserAddressesRequest{
		UserId: uint64(userID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"addresses": resp.Addresses,
	})
}

// EditUserAddress HTTP handler to edit a user address
func (h *UserAddressHandler) EditUserAddress(c *gin.Context) {
	userID := c.GetUint("userid")
	addressStr := c.Param("id")
	addressID,err := strconv.Atoi(addressStr)
	if err!=nil{
		fmt.Println("Failed to convert to uint64")
		return
	}
	var request struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		State      string `json:"state"`
		PostalCode string `json:"postal_code"`
		Country    string `json:"country"`
	}

	// Bind the request JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the EditUserAddress gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userService.EditUserAddress(ctx, &addresspb.EditUserAddressRequest{
		UserId:     uint64(userID),
		AddressId:  uint64(addressID),
		Street:     request.Street,
		City:       request.City,
		State:      request.State,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

// DeleteUserAddress HTTP handler to delete a user address
func (h *UserAddressHandler) DeleteUserAddress(c *gin.Context) {
	userID := c.GetUint("userid")
	var request struct {
		AddressID uint64 `json:"address_id"`
	}

	// Bind the request JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the DeleteUserAddress gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.userService.RemoveUserAddress(ctx, &addresspb.RemoveUserAddressRequest{
		UserId:    uint64(userID),
		AddressId: request.AddressID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}
