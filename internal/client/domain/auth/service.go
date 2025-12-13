package auth

import (
	"context"
	"fmt"

	client "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/auth"
)

type Service interface {
	Register(ctx context.Context, creds Credentials) error
	Login(ctx context.Context, creds Credentials) error
	Logout(ctx context.Context) error
	GetCurrentSession() (Session, error)
}

type service struct {
	client client.AuthClient
	repo   Repository
}

func NewService(client client.AuthClient, repo Repository) Service {
	return &service{
		client: client,
		repo:   repo,
	}
}

func (s *service) Login(ctx context.Context, creds Credentials) error {
	userID, token, err := s.client.Login(ctx, creds.Login, creds.Password)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	err = s.repo.Save(Session{
		User: User{
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

func (s *service) Register(ctx context.Context, creds Credentials) error {
	userID, token, err := s.client.Register(ctx, creds.Login, creds.Password)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	err = s.repo.Save(Session{
		User: User{
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

func (s *service) Logout(ctx context.Context) error {
	err := s.repo.Delete()
	if err != nil {
		return fmt.Errorf("failed to complete logout: %w", err)
	}
	return nil
}

func (s *service) GetCurrentSession() (Session, error) {
	return s.repo.Load()
}
