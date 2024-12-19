package handlers

import (
	"context"

	"github.com/faris-muhammed/e-commerce/admin-service/service"
	adminpb "github.com/faris-muhammed/e-protofiles/adminLogin"
)

type AdminServiceHandler struct {
	adminpb.UnimplementedAdminServiceServer
	adminService service.UserService
}

func NewAdminServiceHandler(adminService service.UserService) *AdminServiceHandler {
	return &AdminServiceHandler{adminService: adminService}
}

func (h *AdminServiceHandler) Login(ctx context.Context, req *adminpb.LoginRequest) (*adminpb.LoginResponse, error) {
	// Perform login and get JWT token
	token, err := h.adminService.Login(req.Username, req.Password)
	if err != nil {
		return &adminpb.LoginResponse{Message: err.Error()}, nil
	}

	// Return JWT token
	return &adminpb.LoginResponse{
		Message: "Login successful",
		Token:   token,
	}, nil
}
