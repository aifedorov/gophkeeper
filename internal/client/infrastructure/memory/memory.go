package memory

import (
	"errors"
	"sync"
)

type Storage struct {
	session *Session
	mu      sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Save(session Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.session = &session

	return nil
}

func (s *Storage) Load() (Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.session == nil {
		return Session{}, errors.New("session not found")
	}

	return *s.session, nil
}

func (s *Storage) Delete() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.session = nil
	return nil
}
