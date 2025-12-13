// Package mw provides gRPC middleware for cross-cutting concerns.
// It includes authentication interceptors for validating JWT tokens and
// injecting user context into request handlers.
package mw

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// jwyKey is the metadata key for extracting JWT access token from gRPC request headers
	jwyKey ContextKey = "access_token"
)

// AuthInterceptor provides JWT-based authentication for gRPC requests.
// It validates tokens from request metadata and injects user context for downstream handlers.
type AuthInterceptor struct {
	jwtSrv  jwt.Service  // JWT token validation service
	authSrv auth.Service // Authentication service for user context management
	logger  *zap.Logger  // Structured logger
}

// NewAuthInterceptor creates a new instance of AuthInterceptor with the provided dependencies.
// It initializes the authentication interceptor that validates JWT tokens in gRPC requests.
func NewAuthInterceptor(jwtSrv jwt.Service, authSrv auth.Service, logger *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSrv:  jwtSrv,
		authSrv: authSrv,
		logger:  logger,
	}
}

// UnaryAuthInterceptor is a gRPC unary interceptor that validates JWT tokens.
// It extracts the access token from request metadata, validates it, extracts the user ID,
// and injects it into the request context for downstream handlers.
// Returns Unauthenticated error if token is missing, invalid, or expired.
func (i *AuthInterceptor) UnaryAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	i.logger.Debug("mw: auth interceptor called", zap.String("method", info.FullMethod))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		i.logger.Debug("mw: no metadata in request")
		return nil, status.Errorf(codes.Unauthenticated, "no metadata")
	}

	token := md.Get(string(jwyKey))
	if len(token) == 0 {
		i.logger.Debug("mw: jwt token not found in metadata")
		return nil, status.Errorf(codes.Unauthenticated, "jwt token not found")
	}

	userID, err := i.jwtSrv.ExtractUserID(token[0])
	if err != nil {
		i.logger.Debug("mw: failed to extract user id from token", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	i.logger.Debug("mw: user authenticated", zap.String("user_id", userID))

	ctx = i.authSrv.SetUserIDToContext(ctx, userID)
	return handler(ctx, req)
}
