package handlers

// import (
// 	"context"

// 	"github.com/faris-muhammed/e-commerce/user-service/service"
// 	userpb "github.com/faris-muhammed/e-protofiles/userlogin"
// )

// type server struct {
// 	userpb.UnimplementedLoginServiceServer
// 	userService service.UserService
// }

// // NewServer creates a new gRPC server instance
// func NewServer(userService service.UserService) *server {
// 	return &server{userService: userService}
// }

// func (s *server) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
// 	resp, err := s.userService.Login(req.Email, req.Password)
// 	if err != nil {
// 		return &userpb.LoginResponse{
// 			Success: false,
// 			Message: resp.Message,
// 		}, err
// 	}

// 	return &userpb.LoginResponse{
// 		Success: resp.Success,
// 		Message: resp.Message,
// 		Token:   resp.Token,
// 	}, nil
// }
