package interseptors

import (
	"context"
	"fmt"

	client "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// jwyKey is the metadata key for extracting JWT access token from gRPC request headers
	jwyKey ContextKey = "access_token"
)

var publicMethods = map[string]bool{
	"/auth.v1.AuthService/Login":    true,
	"/auth.v1.AuthService/Register": true,
}

type AuthInterceptor struct {
	tokenProvider client.TokenProvider
}

func NewAuthInterceptor(tokenProvider client.TokenProvider) *AuthInterceptor {
	return &AuthInterceptor{tokenProvider: tokenProvider}
}

func (i *AuthInterceptor) Interceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if publicMethods[method] {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		accessToken, err := i.tokenProvider.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("authInterseptor: failed to load session: %w", err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, string(jwyKey), accessToken)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
