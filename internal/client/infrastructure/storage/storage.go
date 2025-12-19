package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	filename = "storage.json"
	fileMode = 0600
)

type Storage struct {
	mu sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Save(session Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("storage: failed to marshal session: %w", err)
	}

	err = os.WriteFile(filename, jsonData, fileMode)
	if err != nil {
		return fmt.Errorf("storage: failed to write binary: %w", err)
	}

	return nil
}

func (s *Storage) Load() (Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return Session{}, fmt.Errorf("storage: failed to read binary: %w", err)
	}

	var session Session
	err = json.Unmarshal(jsonData, &session)
	if err != nil {
		return Session{}, fmt.Errorf("storage: failed to unmarshal session: %w", err)
	}

	return session, nil
}

func (s *Storage) Delete() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("storage: failed to remove binary: %w", err)
	}

	return nil
}
