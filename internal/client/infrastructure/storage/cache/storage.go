package cache

import (
	"fmt"
	"os"
	"sync"
)

const (
	filename = "cache.json"
	fileMode = 0600
)

type Storage struct {
	mu sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) SetCredentialVersion(id string, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, err := s.load()
	if err != nil {
		return fmt.Errorf("storage: failed to load secret from cache: %w", err)
	}

	secret.SetCredential(id, version)

	err = s.save(secret)
	if err != nil {
		return fmt.Errorf("storage: failed to save secret to cache: %w", err)
	}
	return nil
}

func (s *Storage) GetCredentialVersion(id string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secret, err := s.load()
	if err != nil {
		return 0, fmt.Errorf("storage: failed to load secret from cache: %w", err)
	}

	val, ok := secret.GetCredentialVersion(id)
	if !ok {
		return 0, fmt.Errorf("storage: credential %s not found in cache", id)
	}
	return val, nil
}

func (s *Storage) DeleteCredentialVersion(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, err := s.load()
	if err != nil {
		return fmt.Errorf("storage: failed to load secret from cache: %w", err)
	}

	secret.DeleteCredential(id)

	err = s.save(secret)
	if err != nil {
		return fmt.Errorf("storage: failed to save secret to cache: %w", err)
	}
	return nil
}

func (s *Storage) SetCardVersion(id string, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, err := s.load()
	if err != nil {
		return fmt.Errorf("storage: failed to load secret from cache: %w", err)
	}

	secret.SetCard(id, version)

	err = s.save(secret)
	if err != nil {
		return fmt.Errorf("storage: failed to save secret to cache: %w", err)
	}
	return nil
}

func (s *Storage) GetCardVersion(id string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secret, err := s.load()
	if err != nil {
		return 0, fmt.Errorf("storage: failed to load secret from cache: %w", err)
	}

	val, ok := secret.GetCardVersion(id)
	if !ok {
		return 0, fmt.Errorf("storage: card %s not found in cache", id)
	}
	return val, nil
}

func (s *Storage) SetFileVersion(id string, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, err := s.load()
	if err != nil {
		return fmt.Errorf("storage: failed to load secret from cache: %w", err)
	}

	secret.SetFileVersion(id, version)

	err = s.save(secret)
	if err != nil {
		return fmt.Errorf("storage: failed to save secret to cache: %w", err)
	}
	return nil
}

func (s *Storage) GetFileVersion(id string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secret, err := s.load()
	if err != nil {
		return 0, fmt.Errorf("storage: failed to load secret from cache: %w", err)
	}

	val, ok := secret.GetFileVersion(id)
	if !ok {
		return 0, fmt.Errorf("storage: file %s not found in cache", id)
	}
	return val, nil
}

func (s *Storage) ClearAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage: failed to remove cache file: %w", err)
	}

	return nil
}
