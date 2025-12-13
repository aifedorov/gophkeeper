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
	// jwyKey is the context key for storing JWT token
	jwyKey ContextKey = "access_token"
)

type AuthInterceptor struct {
	jwtSrv  jwt.Service
	authSrv auth.Service
	logger  *zap.Logger
}

func NewAuthInterceptor(jwtSrv jwt.Service, authSrv auth.Service, logger *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSrv:  jwtSrv,
		authSrv: authSrv,
		logger:  logger,
	}
}

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
