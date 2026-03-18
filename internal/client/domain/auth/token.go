// Package auth provides session provider implementation for the GophKeeper client.
package auth

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
)

// sessionProvider implements the SessionProvider interface using local storage.
type sessionProvider struct {
	sessionStore interfaces.SessionStore
}

// NewSessionProvider creates a new session provider that retrieves sessions from local storage.
func NewSessionProvider(sessionStore interfaces.SessionStore) interfaces.SessionProvider {
	return &sessionProvider{sessionStore: sessionStore}
}

// GetSession retrieves the current user session from local storage.
// Returns an error if no session is found or if session loading fails.
func (p *sessionProvider) GetSession(_ context.Context) (shared.Session, error) {
	session, err := p.sessionStore.Load()
	if err != nil {
		return shared.Session{}, ErrSessionNotFound
	}
	return session, nil
}
