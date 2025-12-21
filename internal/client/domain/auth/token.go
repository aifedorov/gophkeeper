// Package auth provides session provider implementation for the GophKeeper client.
package auth

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage"
)

// sessionProvider implements the SessionProvider interface using local storage.
type sessionProvider struct {
	sessionStore *storage.Storage
}

// NewSessionProvider creates a new session provider that retrieves sessions from local storage.
func NewSessionProvider(sessionStore *storage.Storage) interfaces.SessionProvider {
	return &sessionProvider{sessionStore: sessionStore}
}

// GetSession retrieves the current user session from local storage.
// Returns an error if no session is found or if session loading fails.
func (p *sessionProvider) GetSession(_ context.Context) (interfaces.Session, error) {
	session, err := p.sessionStore.Load()
	if err != nil {
		return interfaces.Session{}, ErrSessionNotFound
	}
	return interfaces.NewSession(
		session.GetAccessToken(),
		session.GetEncryptionKey(),
		session.GetUserID(),
	), nil
}
