package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/yowie645/sso-grpc-go/internal/app/grpc"
	"github.com/yowie645/sso-grpc-go/internal/services/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	authService := auth.New(log, storagePath, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
