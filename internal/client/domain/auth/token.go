package auth

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage"
)

type sessionProvider struct {
	sessionStore *storage.Storage
}

func NewSessionProvider(sessionStore *storage.Storage) interfaces.SessionProvider {
	return &sessionProvider{sessionStore: sessionStore}
}

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
