package auth

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	client "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/auth"
)

type Service interface {
	Register(ctx context.Context, creds interfaces.Credentials) error
	Login(ctx context.Context, creds interfaces.Credentials) error
	Logout(ctx context.Context) error
	GetCurrentSession() (interfaces.Session, error)
}

type service struct {
	client client.AuthClient
	repo   interfaces.Repository
}

func NewService(client client.AuthClient, repo interfaces.Repository) Service {
	return &service{
		client: client,
		repo:   repo,
	}
}

func (s *service) Login(ctx context.Context, creds interfaces.Credentials) error {
	userID, token, err := s.client.Login(ctx, creds.Login, creds.Password)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	err = s.repo.Save(interfaces.Session{
		User: interfaces.User{
			ID:    userID,
			Login: creds.Login,
		},
		AccessToken: token,
	})
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (s *service) Register(ctx context.Context, creds interfaces.Credentials) error {
	userID, token, err := s.client.Register(ctx, creds.Login, creds.Password)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	err = s.repo.Save(interfaces.Session{
		User: interfaces.User{
			ID:    userID,
			Login: creds.Login,
		},
		AccessToken: token,
	})
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (s *service) Logout(_ context.Context) error {
	err := s.repo.Delete()
	if err != nil {
		return fmt.Errorf("failed to complete logout: %w", err)
	}
	return nil
}

func (s *service) GetCurrentSession() (interfaces.Session, error) {
	return s.repo.Load()
}
