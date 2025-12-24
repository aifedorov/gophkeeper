package cache

import (
	"encoding/json"
	"fmt"
	"os"
)

func (s *Storage) load() (*Secret, error) {
	jsonData, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return NewSecret(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("storage: failed to read cache: %w", err)
	}

	var secret Secret
	err = json.Unmarshal(jsonData, &secret)
	if err != nil {
		return nil, fmt.Errorf("storage: failed to unmarshal cache: %w", err)
	}

	if secret.Credentials == nil {
		secret.Credentials = make(map[string]int64)
	}
	if secret.Cards == nil {
		secret.Cards = make(map[string]int64)
	}
	if secret.Files == nil {
		secret.Files = make(map[string]int64)
	}

	return &secret, nil
}

func (s *Storage) save(secret *Secret) error {
	jsonData, err := json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("storage: failed to marshal secret: %w", err)
	}

	err = os.WriteFile(filename, jsonData, fileMode)
	if err != nil {
		return fmt.Errorf("storage: failed to write secret to cache: %w", err)
	}

	return nil
}
