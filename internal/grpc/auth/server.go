package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	authv1 "github.com/yowie645/protos-sso-grcp-go/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	validate *validator.Validate
	auth     Auth
}

func RegisterServer(gRPC *grpc.Server, auth Auth) {
	api := &serverAPI{
		validate: validator.New(),
		auth:     auth,
	}
	authv1.RegisterAuthServer(gRPC, api)
}

func (s *serverAPI) validateRequest(req interface{}) error {
	if err := s.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			return status.Errorf(codes.InvalidArgument, "invalid field %s: %v", e.Field(), e.Tag())
		}
	}
	return nil
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}

	return &authv1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "admin check failed: %v", err)
	}

	return &authv1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
