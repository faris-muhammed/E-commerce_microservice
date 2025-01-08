package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/faris-muhammed/e-commerce/user-service/models"
	"github.com/faris-muhammed/e-commerce/user-service/repository"
	addresspb "github.com/faris-muhammed/e-protofiles/address"
	"gorm.io/gorm"
)

// UserAddressServiceServer implements the UserService gRPC service
type UserAddressServiceServer struct {
	addresspb.UnimplementedUserServiceServer
	userRepo repository.UserAddressRepository
}

func NewUserAddressService(repo repository.UserAddressRepository) *UserAddressServiceServer {
	return &UserAddressServiceServer{userRepo: repo}
}

// AddUserAddress adds an address for a user
func (s *UserAddressServiceServer) AddUserAddress(ctx context.Context, req *addresspb.AddUserAddressRequest) (*addresspb.AddUserAddressResponse, error) {
	// Validate request data
	if req.UserId == 0 || req.Street == "" || req.City == "" || req.State == "" || req.PostalCode == "" || req.Country == "" {
		return nil, errors.New("all fields are required")
	}

	// Check if the address already exists
	existingAddress, err := s.userRepo.GetAddressByStreet(req.Street)
	if err != nil && err != gorm.ErrRecordNotFound { // Ensure it's not a "not found" error
		return nil, fmt.Errorf("failed to check address existence: %w", err)
	}
	if existingAddress != nil && existingAddress.UserID == req.UserId {
		return nil, errors.New("address already exists for this user")
	}

	// Create a new address in the database
	address := models.Address{
		UserID:     req.UserId,
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
	}

	err = s.userRepo.CreateUserAddress(&address)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	// Return the address ID and a success message
	return &addresspb.AddUserAddressResponse{
		Message:   "Address added successfully",
		AddressId: address.ID,
	}, nil
}

// GetUserAddresses retrieves all addresses of a user
func (s *UserAddressServiceServer) GetUserAddresses(ctx context.Context, req *addresspb.GetUserAddressesRequest) (*addresspb.GetUserAddressesResponse, error) {
	// Fetch addresses for the user from the repository
	addresses, err := s.userRepo.GetUserAddresses(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}

	// Convert addresses to gRPC format
	var grpcAddresses []*addresspb.Address
	for _, address := range addresses {
		grpcAddresses = append(grpcAddresses, &addresspb.Address{
			AddressId:  address.ID,
			Street:     address.Street,
			City:       address.City,
			State:      address.State,
			PostalCode: address.PostalCode,
			Country:    address.Country,
		})
	}

	// Return the addresses
	return &addresspb.GetUserAddressesResponse{
		Addresses: grpcAddresses,
	}, nil
}

// EditUserAddress updates an existing user address
func (s *UserAddressServiceServer) EditUserAddress(ctx context.Context, req *addresspb.EditUserAddressRequest) (*addresspb.EditUserAddressResponse, error) {

	// Check if the address exists
	existingAddress, err := s.userRepo.GetAddressByID(req.AddressId)
	if err != nil {
		return nil, fmt.Errorf("failed to find address: %w", err)
	}

	// Update the address details
	if req.Street != "" {
		existingAddress.Street = req.Street
	}
	if req.City != "" {
		existingAddress.City = req.City
	}
	if req.State != "" {
		existingAddress.State = req.State
	}
	if req.PostalCode != "" {
		existingAddress.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		existingAddress.Country = req.Country
	}
	// Save the updated address in the repository
	_, erro := s.userRepo.UpdateUserAddress(existingAddress)
	if erro != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	// Return a success message
	return &addresspb.EditUserAddressResponse{
		Message: "Address updated successfully",
	}, nil
}

// RemoveUserAddress removes a user address
func (s *UserAddressServiceServer) RemoveUserAddress(ctx context.Context, req *addresspb.RemoveUserAddressRequest) (*addresspb.RemoveUserAddressResponse, error) {
	// Validate request
	if req.AddressId == 0 || req.UserId == 0 {
		return nil, errors.New("address ID and user ID are required")
	}

	// Check if the address exists
	address, err := s.userRepo.GetAddressByID(req.AddressId)
	if err != nil {
		return nil, fmt.Errorf("failed to find address: %w", err)
	}

	// Remove the address from the repository
	err = s.userRepo.DeleteUserAddress(address)
	if err != nil {
		return nil, fmt.Errorf("failed to delete address: %w", err)
	}

	// Return a success message
	return &addresspb.RemoveUserAddressResponse{
		Message: "Address removed successfully",
	}, nil
}
