package auth

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage"
)

type tokenProvider struct {
	sessionStore *storage.Storage
}

func NewTokeProvider(sessionStore *storage.Storage) interfaces.TokenProvider {
	return &tokenProvider{sessionStore: sessionStore}
}

func (p *tokenProvider) GetToken(_ context.Context) (string, error) {
	session, err := p.sessionStore.Load()
	if err != nil {
		return "", fmt.Errorf("tokenProvider: failed to get token: %w", err)
	}
	return session.AccessToken, nil
}
