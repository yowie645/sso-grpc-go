package auth

import (
	"context"

	authv1 "github.com/yowie645/protos-sso-grcp-go/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

func RegisterServer(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	panic("implement me")
}
