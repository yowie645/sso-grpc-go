package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yowie645/sso-grpc-go/internal/domain/models"
)

func NewToken(user models.User, app models.App, ttl time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":    user.ID,
		"email":  user.Email,
		"app_id": app.ID,
		"exp":    time.Now().Add(ttl).Unix(),
	})
	return token.SignedString([]byte(secret))
}
