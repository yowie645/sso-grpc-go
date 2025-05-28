package auth

import (
	"context"
	"log/slog"
	"time"
)

// Auth - (пока заглушка)
type Auth struct {
	log         *slog.Logger
	storagePath string
	tokenTTL    time.Duration
}

// IsAdmin implements auth.Auth.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("unimplemented")
}

// Login implements auth.Auth.
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (token string, err error) {
	panic("unimplemented")
}

// RegisterNewUser implements auth.Auth.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error) {
	panic("unimplemented")
}

// New - (пока заглушка)
func New(log *slog.Logger, storagePath string, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		storagePath: storagePath,
		tokenTTL:    tokenTTL,
	}
}
