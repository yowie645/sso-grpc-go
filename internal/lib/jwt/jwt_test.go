package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yowie645/sso-grpc-go/internal/domain/models"
)

const (
	testSecret = "test-secret-key-1234567890" // Тестовый секретный ключ
)

func TestNewToken(t *testing.T) {
	now := time.Now()
	testUser := models.User{
		ID:    1,
		Email: "test@example.com",
	}
	testApp := models.App{
		ID: 1,
	}
	ttl := 1 * time.Hour

	verifyToken := func(t *testing.T, tokenStr string) jwt.MapClaims {
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(testSecret), nil
		})
		require.NoError(t, err)
		require.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)
		return claims
	}

	t.Run("successful token generation", func(t *testing.T) {
		tokenStr, err := NewToken(testUser, testApp, ttl, testSecret)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenStr)

		claims := verifyToken(t, tokenStr)

		assert.Equal(t, float64(testUser.ID), claims["uid"])
		assert.Equal(t, testUser.Email, claims["email"])
		assert.Equal(t, float64(testApp.ID), claims["app_id"])

		exp, ok := claims["exp"].(float64)
		require.True(t, ok)
		assert.InDelta(t, now.Add(ttl).Unix(), int64(exp), 1)
	})

	t.Run("different users and apps", func(t *testing.T) {
		testCases := []struct {
			name  string
			user  models.User
			app   models.App
			check func(t *testing.T, claims jwt.MapClaims)
		}{
			{
				name: "zero values",
				user: models.User{},
				app:  models.App{},
				check: func(t *testing.T, claims jwt.MapClaims) {
					assert.Equal(t, float64(0), claims["uid"])
					assert.Equal(t, "", claims["email"])
					assert.Equal(t, float64(0), claims["app_id"])
				},
			},
			{
				name: "max int64 user id",
				user: models.User{ID: 1<<63 - 1},
				app:  models.App{ID: 1},
				check: func(t *testing.T, claims jwt.MapClaims) {
					assert.Equal(t, float64(1<<63-1), claims["uid"])
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tokenStr, err := NewToken(tc.user, tc.app, ttl, testSecret)
				require.NoError(t, err)

				claims := verifyToken(t, tokenStr)
				tc.check(t, claims)
			})
		}
	})

	t.Run("invalid signing method", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodNone)
		tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		_, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(testSecret), nil
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'none' signature type is not allowed")
	})
}

func TestTokenExpiration(t *testing.T) {
	testUser := models.User{ID: 1, Email: "test@example.com"}
	testApp := models.App{ID: 1}

	verifyToken := func(_ *testing.T, tokenStr string) (jwt.MapClaims, error) {
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(testSecret), nil
		})
		if err != nil {
			return nil, err
		}
		return token.Claims.(jwt.MapClaims), nil
	}

	t.Run("zero duration", func(t *testing.T) {
		tokenStr, err := NewToken(testUser, testApp, 0, testSecret)
		require.NoError(t, err)

		parser := new(jwt.Parser)
		token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
		require.NoError(t, err)

		claims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)

		exp, ok := claims["exp"].(float64)
		require.True(t, ok)
		assert.InDelta(t, time.Now().Unix(), int64(exp), 1)

		_, err = verifyToken(t, tokenStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("negative duration", func(t *testing.T) {
		tokenStr, err := NewToken(testUser, testApp, -1*time.Hour, testSecret)
		require.NoError(t, err)

		parser := new(jwt.Parser)
		token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
		require.NoError(t, err)

		claims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)

		exp, ok := claims["exp"].(float64)
		require.True(t, ok)
		assert.Less(t, int64(exp), time.Now().Unix())

		_, err = verifyToken(t, tokenStr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})
}
