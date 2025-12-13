package grpc

import (
	"context"
	"fmt"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/auth/v1"
	"google.golang.org/grpc"
)

type AuthClient interface {
	Register(ctx context.Context, login, pass string) (userID, token string, err error)
	Login(ctx context.Context, login, pass string) (userID, token string, err error)
}

type authClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) AuthClient {
	return &authClient{
		client: pb.NewAuthServiceClient(conn),
	}
}

func (c *authClient) Register(ctx context.Context, login, pass string) (userID, token string, err error) {
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Login:    &login,
		Password: &pass,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to register: %w", err)
	}
	return resp.GetUserId(), resp.GetAccessToken(), nil
}

func (c *authClient) Login(ctx context.Context, login, pass string) (userID, token string, err error) {
	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Login:    &login,
		Password: &pass,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to login: %w", err)
	}
	return resp.GetUserId(), resp.GetAccessToken(), nil
}
