// Package auth provides authentication services for the GophKeeper client.
//
// This package implements the client-side authentication logic including user registration,
// login, logout, and session management. It communicates with the server via gRPC and
// manages local session storage.
package auth

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/client"
)

// Service defines the interface for client-side authentication operations.
type Service interface {
	// Register creates a new user account with the provided credentials.
	// It sends a registration request to the server, receives a session (JWT token and encryption key),
	// and saves it to local storage. Returns an error if registration fails or session save fails.
	Register(ctx context.Context, creds interfaces.Credentials) error
	// Login authenticates a user with the provided credentials.
	// It sends a login request to the server, receives a session (JWT token and encryption key),
	// and saves it to local storage. Returns an error if authentication fails or session save fails.
	Login(ctx context.Context, creds interfaces.Credentials) error
	// Logout removes the current user session from local storage.
	// Returns an error if session deletion fails.
	Logout(ctx context.Context) error
	// GetCurrentSession retrieves the current user session from local storage.
	// Returns an error if no session is found or if session loading fails.
	GetCurrentSession() (interfaces.Session, error)
}

// service implements the Service interface for client-side authentication.
type service struct {
	client client.AuthClient
	repo   interfaces.Repository
}

// NewService creates a new instance of the authentication service with the provided dependencies.
// It initializes the service with a gRPC auth client and a local session repository.
func NewService(client client.AuthClient, repo interfaces.Repository) Service {
	return &service{
		client: client,
		repo:   repo,
	}
}

// Login authenticates a user with the provided credentials.
// It sends a login request to the server, receives a session (JWT token and encryption key),
// and saves it to local storage. Returns an error if authentication fails or session save fails.
func (s *service) Login(ctx context.Context, creds interfaces.Credentials) error {
	session, err := s.client.Login(ctx, creds.GetLogin(), creds.GetPassword())
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	err = s.repo.Save(session)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// Register creates a new user account with the provided credentials.
// It sends a registration request to the server, receives a session (JWT token and encryption key),
// and saves it to local storage. Returns an error if registration fails or session save fails.
func (s *service) Register(ctx context.Context, creds interfaces.Credentials) error {
	session, err := s.client.Register(ctx, creds.GetLogin(), creds.GetPassword())
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	err = s.repo.Save(session)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// Logout removes the current user session from local storage.
// Returns an error if session deletion fails.
func (s *service) Logout(_ context.Context) error {
	err := s.repo.Delete()
	if err != nil {
		return fmt.Errorf("failed to complete logout: %w", err)
	}
	return nil
}

// GetCurrentSession retrieves the current user session from local storage.
// Returns an error if no session is found or if session loading fails.
func (s *service) GetCurrentSession() (interfaces.Session, error) {
	return s.repo.Load()
}
