package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/auth/v1"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=auth.go -destination=mock_auth_client.go -package=client

type AuthClient interface {
	Register(ctx context.Context, login, pass string) (session shared.Session, err error)
	Login(ctx context.Context, login, pass string) (session shared.Session, err error)
}

type authClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) AuthClient {
	return &authClient{
		client: pb.NewAuthServiceClient(conn),
	}
}

func (c *authClient) Register(ctx context.Context, login, pass string) (session shared.Session, err error) {
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Login:    &login,
		Password: &pass,
	})
	if err != nil {
		return shared.Session{}, fmt.Errorf("authClient: failed to register: %w", err)
	}

	return shared.NewSession(
		resp.GetAccessToken(),
		base64.StdEncoding.EncodeToString(resp.GetEncryptionKey()),
		resp.GetUserId(),
		resp.GetLogin(),
	), nil
}

func (c *authClient) Login(ctx context.Context, login, pass string) (session shared.Session, err error) {
	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Login:    &login,
		Password: &pass,
	})
	if err != nil {
		return shared.Session{}, fmt.Errorf("authClient: failed to login: %w", err)
	}

	return shared.NewSession(
		resp.GetAccessToken(),
		base64.StdEncoding.EncodeToString(resp.GetEncryptionKey()),
		resp.GetUserId(),
		resp.GetLogin(),
	), nil
}
