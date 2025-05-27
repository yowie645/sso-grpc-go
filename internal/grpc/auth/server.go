package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	authv1 "github.com/yowie645/protos-sso-grcp-go/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	validate *validator.Validate
}

func RegisterServer(gRPC *grpc.Server) {
	api := &serverAPI{
		validate: validator.New(),
	}
	authv1.RegisterAuthServer(gRPC, api)
}

func (s *serverAPI) validateRequest(req interface{}) error {
	if err := s.validate.Struct(req); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}
	return nil
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	return &authv1.LoginResponse{
		Token: req.GetEmail(),
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	if len(req.GetPassword()) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	panic("implement me")
}
