package handlers

import (
	"context"

	"github.com/faris-muhammed/e-commerce/seller-service/service"
	sellerpb "github.com/faris-muhammed/e-protofiles/sellerlogin"
)

type SellerServiceHandler struct {
	sellerpb.UnimplementedSellerServiceServer
	sellerService service.UserService
}

func NewSellerServiceHandler(sellerService service.UserService) *SellerServiceHandler {
	return &SellerServiceHandler{sellerService: sellerService}
}

func (h *SellerServiceHandler) Login(ctx context.Context, req *sellerpb.LoginRequest) (*sellerpb.LoginResponse, error) {
	// Perform login and get JWT token
	token, err := h.sellerService.Login(req.Username, req.Password)
	if err != nil {
		return &sellerpb.LoginResponse{Message: err.Error()}, nil
	}

	// Return JWT token
	return &sellerpb.LoginResponse{
		Message: "Login successful",
		Token:   token,
	}, nil
}
