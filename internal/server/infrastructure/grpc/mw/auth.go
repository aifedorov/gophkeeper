package mw

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/jwt"
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

type Interceptor struct {
	jwtSrv  jwt.Service
	authSrv auth.Service
}

func (i *Interceptor) UnaryAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata")
	}

	token := md.Get(string(jwyKey))
	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "jwt token not found")
	}

	userID, err := i.jwtSrv.ExtractUserID(token[0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	ctx = i.authSrv.SetUserIDToContext(ctx, userID)
	return handler(ctx, req)
}
