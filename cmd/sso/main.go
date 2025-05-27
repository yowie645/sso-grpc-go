package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yowie645/sso-grpc-go/internal/app"
	"github.com/yowie645/sso-grpc-go/internal/config"
	"github.com/yowie645/sso-grpc-go/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := setupLogger(cfg.Env)
	log.Info("start app", slog.Any("cfg", cfg))

	// TODO инициализация логгера

	aplication := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	aplication.GRPCSrv.MustRun()
	// TODO запуск приложения

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Info("stop signal received", slog.String("signal", sign.String()))

	aplication.GRPCSrv.Stop()
	log.Info("stop app")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug}}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
