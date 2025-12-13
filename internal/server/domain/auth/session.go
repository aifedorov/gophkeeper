package auth

import (
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SessionStore interface {
	GetEncryptionKey(userID uuid.UUID) ([]byte, bool)
	Set(userID uuid.UUID, key []byte)
}

type sessionStore struct {
	keys   sync.Map
	logger *zap.Logger
}

func NewSessionStore(logger *zap.Logger) SessionStore {
	return &sessionStore{
		keys:   sync.Map{},
		logger: logger,
	}
}

func (s *sessionStore) GetEncryptionKey(userID uuid.UUID) ([]byte, bool) {
	s.logger.Debug("session: retrieving encryption key", zap.String("user_id", userID.String()))
	if val, ok := s.keys.Load(userID); ok {
		s.logger.Debug("session: encryption key found", zap.String("user_id", userID.String()))
		return val.([]byte), true
	}
	s.logger.Debug("session: encryption key not found", zap.String("user_id", userID.String()))
	return nil, false
}

func (s *sessionStore) Set(userID uuid.UUID, key []byte) {
	s.logger.Debug("session: storing encryption key", zap.String("user_id", userID.String()))
	s.keys.Store(userID, key)
}
