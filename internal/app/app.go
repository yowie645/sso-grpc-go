package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/yowie645/sso-grpc-go/internal/app/grpc"
	"github.com/yowie645/sso-grpc-go/internal/services/auth"
	"github.com/yowie645/sso-grpc-go/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	jwtSecret string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(
		log,
		storage,
		storage,
		storage,
		storagePath,
		tokenTTL,
		jwtSecret,
	)
	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
