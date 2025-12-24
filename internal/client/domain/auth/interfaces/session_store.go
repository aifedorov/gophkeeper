// Package auth provides authentication services for the GophKeeper client.
package interfaces

import "github.com/aifedorov/gophkeeper/internal/client/domain/shared"

//go:generate mockgen -source=session_store.go -destination=mock_session_store_test.go -package=auth

// SessionStore defines the interface for local session storage operations.
// Implementations typically store sessions in local files or secure key-value stores.
type SessionStore interface {
	// Save stores a session in local storage.
	// Returns an error if session storage fails.
	Save(session shared.Session) error
	// Load retrieves a session from local storage.
	// Returns an error if no session is found or if session loading fails.
	Load() (shared.Session, error)
	// Delete removes the current session from local storage.
	// Returns an error if session deletion fails.
	Delete() error
}
