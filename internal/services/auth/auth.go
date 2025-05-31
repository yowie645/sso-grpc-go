package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/yowie645/sso-grpc-go/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

// Auth - (пока заглушка)
type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	storagePath string
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (userID int64, err error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	storagePath string,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		storagePath: storagePath,
		tokenTTL:    tokenTTL,
	}
}

// IsAdmin implements auth.Auth.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("unimplemented")
}

// Login implements auth.Auth.
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	panic("unimplemented")
}

// RegisterNewUser implements auth.Auth.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("register new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("registered new user")

	return id, nil
}
