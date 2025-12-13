package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	testSecretKey = "test-secret-key-for-jwt-signing"
	testTokenExp  = 24 * time.Hour
)

func TestNewService(t *testing.T) {
	logger := zap.NewNop()
	service := NewService(testSecretKey, testTokenExp, logger)

	require.NotNil(t, service)
}

func TestService_IssueToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "valid user ID",
			userID:  uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: false,
		},
		{
			name:    "numeric user ID",
			userID:  "12345",
			wantErr: false,
		},
		{
			name:    "user ID with special characters",
			userID:  "user@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := NewService(testSecretKey, testTokenExp, zap.NewNop())
			token, err := service.IssueToken(tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestService_ExtractUserID(t *testing.T) {
	t.Parallel()

	service := NewService(testSecretKey, testTokenExp, zap.NewNop())

	t.Run("extract from valid token", func(t *testing.T) {
		userID := uuid.New().String()

		token, err := service.IssueToken(userID)
		require.NoError(t, err)

		extractedUserID, err := service.ExtractUserID(token)
		require.NoError(t, err)
		assert.Equal(t, userID, extractedUserID)
	})

	t.Run("empty token string", func(t *testing.T) {
		_, err := service.ExtractUserID("")
		assert.ErrorIs(t, err, ErrEmptyToken)
	})

	t.Run("invalid token format", func(t *testing.T) {
		_, err := service.ExtractUserID("invalid-token")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("token with wrong secret", func(t *testing.T) {
		wrongService := NewService("wrong-secret", testTokenExp, zap.NewNop())
		userID := uuid.New().String()

		token, err := wrongService.IssueToken(userID)
		require.NoError(t, err)

		_, err = service.ExtractUserID(token)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("expired token", func(t *testing.T) {
		expiredService := NewService(testSecretKey, -1*time.Hour, zap.NewNop())
		userID := uuid.New().String()

		token, err := expiredService.IssueToken(userID)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		_, err = service.ExtractUserID(token)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("token with invalid signing method", func(t *testing.T) {
		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(testTokenExp)),
			},
			UserID: uuid.New().String(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		_, err = service.ExtractUserID(tokenString)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("malformed token", func(t *testing.T) {
		malformedTokens := []string{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",                             // Only header
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0", // No signature
			"not.a.token", // Invalid format
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid-payload.signature", // Invalid payload
		}

		for _, token := range malformedTokens {
			_, err := service.ExtractUserID(token)
			assert.ErrorIs(t, err, ErrInvalidToken, "token: %s", token)
		}
	})
}

func TestService_IssueAndExtract_Integration(t *testing.T) {
	t.Parallel()

	service := NewService(testSecretKey, testTokenExp, zap.NewNop())

	testCases := []struct {
		name   string
		userID string
	}{
		{
			name:   "UUID user ID",
			userID: uuid.New().String(),
		},
		{
			name:   "numeric user ID",
			userID: "123456",
		},
		{
			name:   "email-like user ID",
			userID: "user@example.com",
		},
		{
			name:   "empty user ID",
			userID: "",
		},
		{
			name:   "long user ID",
			userID: "very-long-user-id-with-many-characters-" + uuid.New().String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			token, err := service.IssueToken(tc.userID)
			require.NoError(t, err)
			assert.NotEmpty(t, token)

			extractedUserID, err := service.ExtractUserID(token)
			require.NoError(t, err)
			assert.Equal(t, tc.userID, extractedUserID)
		})
	}
}

func TestService_TokenExpiration(t *testing.T) {
	t.Parallel()

	t.Run("token valid before expiration", func(t *testing.T) {
		longExp := 1 * time.Hour
		service := NewService(testSecretKey, longExp, zap.NewNop())
		userID := uuid.New().String()

		token, err := service.IssueToken(userID)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		extractedUserID, err := service.ExtractUserID(token)
		require.NoError(t, err)
		assert.Equal(t, userID, extractedUserID)
	})
}

func TestService_DifferentSecretKeys(t *testing.T) {
	t.Parallel()

	service1 := NewService("secret-key-1", testTokenExp, zap.NewNop())
	service2 := NewService("secret-key-2", testTokenExp, zap.NewNop())
	userID := uuid.New().String()

	token, err := service1.IssueToken(userID)
	require.NoError(t, err)

	_, err = service2.ExtractUserID(token)
	assert.ErrorIs(t, err, ErrInvalidToken)

	extractedUserID, err := service1.ExtractUserID(token)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}

func TestService_MultipleTokens(t *testing.T) {
	t.Parallel()

	service := NewService(testSecretKey, testTokenExp, zap.NewNop())

	users := make(map[string]string)
	for i := 0; i < 10; i++ {
		userID := uuid.New().String()
		token, err := service.IssueToken(userID)
		require.NoError(t, err)
		users[token] = userID
	}

	for token, expectedUserID := range users {
		extractedUserID, err := service.ExtractUserID(token)
		require.NoError(t, err)
		assert.Equal(t, expectedUserID, extractedUserID)
	}
}

func TestClaims(t *testing.T) {
	t.Parallel()

	t.Run("claims structure", func(t *testing.T) {
		userID := uuid.New().String()
		expiresAt := time.Now().Add(1 * time.Hour)

		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
			},
			UserID: userID,
		}

		assert.Equal(t, userID, claims.UserID)
		assert.NotNil(t, claims.ExpiresAt)
	})
}

func TestErrorTypes(t *testing.T) {
	t.Parallel()

	t.Run("error constants are defined", func(t *testing.T) {
		assert.NotNil(t, ErrEmptyToken)
		assert.NotNil(t, ErrInvalidToken)
		assert.NotNil(t, ErrInvalidSigningMethod)

		assert.Equal(t, "token is empty", ErrEmptyToken.Error())
		assert.Equal(t, "invalid token", ErrInvalidToken.Error())
		assert.Equal(t, "unexpected signing method", ErrInvalidSigningMethod.Error())
	})
}
