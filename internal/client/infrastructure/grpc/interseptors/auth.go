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
	// accessTokenKey is the metadata key for extracting JWT access token from gRPC request headers
	accessTokenKey ContextKey = "access_token"
	// encryptionKeyKey is the metadata key for extracting encryption key from gRPC request headers
	encryptionKeyKey ContextKey = "encryption_key"
)

var publicMethods = map[string]bool{
	"/auth.v1.AuthService/Login":    true,
	"/auth.v1.AuthService/Register": true,
}

type AuthInterceptor struct {
	sessionProvider client.SessionProvider
}

func NewAuthInterceptor(tokenProvider client.SessionProvider) *AuthInterceptor {
	return &AuthInterceptor{sessionProvider: tokenProvider}
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

		session, err := i.sessionProvider.GetSession(ctx)
		if err != nil {
			fmt.Printf("method: %s, err: %v\n", method, err.Error())
			return fmt.Errorf("authInterseptor: failed to load session: %w", err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx,
			string(accessTokenKey), session.GetAccessToken(),
			string(encryptionKeyKey), session.GetEncryptionKey(),
		)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if publicMethods[method] {
			return streamer(ctx, desc, cc, method, opts...)
		}

		session, err := i.sessionProvider.GetSession(ctx)
		if err != nil {
			fmt.Printf("method: %s, err: %v\n", method, err.Error())
			return nil, fmt.Errorf("authInterseptor: failed to load session: %w", err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx,
			string(accessTokenKey), session.GetAccessToken(),
			string(encryptionKeyKey), session.GetEncryptionKey(),
		)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
