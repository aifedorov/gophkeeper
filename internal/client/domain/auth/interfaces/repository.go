// Package interfaces provides session repository interface for the GophKeeper client.
package interfaces

//go:generate mockgen -source=interfaces.go -destination=mock_repository_test.go -package=auth

// Repository defines the interface for local session storage operations.
// Implementations typically store sessions in local files or key-value stores.
type Repository interface {
	// Save stores a session in local storage.
	// Returns an error if session storage fails.
	Save(session Session) error
	// Load retrieves a session from local storage.
	// Returns an error if no session is found or if session loading fails.
	Load() (Session, error)
	// Delete removes the current session from local storage.
	// Returns an error if session deletion fails.
	Delete() error
}
