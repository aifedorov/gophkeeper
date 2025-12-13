package interfaces

import "github.com/google/uuid"

//go:generate mockgen -source=session.go -destination=mocks/mock_session.go -package=mocks

// SessionStore defines the interface for session storage operations.
type SessionStore interface {
	// GetEncryptionKey retrieves the encryption key for a user.
	// Returns the key and true if found, nil and false otherwise.
	GetEncryptionKey(userID uuid.UUID) ([]byte, bool)
	// Set stores the encryption key for a user.
	Set(userID uuid.UUID, key []byte)
}
