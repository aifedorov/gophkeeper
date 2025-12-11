package auth

import (
	"sync"

	"github.com/google/uuid"
)

type SessionStore struct {
	keys sync.Map
}

func NewSessionStore() *SessionStore {
	return &SessionStore{keys: sync.Map{}}
}

func (s *SessionStore) GetEncryptionKey(userID uuid.UUID) ([]byte, bool) {
	if val, ok := s.keys.Load(userID); ok {
		return val.([]byte), true
	}
	return nil, false
}

func (s *SessionStore) Set(userID uuid.UUID, key []byte) {
	s.keys.Store(userID, key)
}
