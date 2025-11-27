package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"errors"
	"time"
)

var (
	// ErrEmptyToken is returned when an empty token string is provided for parsing.
	ErrEmptyToken = errors.New("token is empty")

	// ErrInvalidToken is returned when token parsing or validation fails.
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidSigningMethod is returned when the token uses an unexpected signing method.
	ErrInvalidSigningMethod = errors.New("unexpected signing method")
)

type Service interface {
	IssueToken(userID string) (string, error)
	ValidateToken(tokenString string) (string, error)
}

type service struct {
	secretKey string
	tokenExp  time.Duration
	logger    *zap.Logger
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func New(secretKey string, tokenExp time.Duration) Service {
	return &service{
		secretKey: secretKey,
		tokenExp:  tokenExp,
	}
}

func (s *service) IssueToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		s.logger.Error("jwt: failed to sign token", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (s *service) ValidateToken(tokenString string) (string, error) {
	if tokenString == "" {
		s.logger.Error("jwt: empty token")
		return "", ErrEmptyToken
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", ErrInvalidSigningMethod, t.Header["alg"])
			}
			return []byte(s.secretKey), nil
		})

	if err != nil {
		s.logger.Error("jwt: error parsing token", zap.Error(err))
		return "", ErrInvalidToken
	}

	if !token.Valid {
		s.logger.Error("jwt: invalid token")
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}
