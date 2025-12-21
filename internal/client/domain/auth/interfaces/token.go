// Package interfaces provides session provider interface for the GophKeeper client.
package interfaces

import "context"

//go:generate mockgen -source=token.go -destination=mock_session_provider_test.go -package=interfaces

// SessionProvider defines the interface for retrieving user sessions.
// Implementations typically load sessions from local storage.
type SessionProvider interface {
	// GetSession retrieves the current user session.
	// Returns an error if no session is found or if session loading fails.
	GetSession(ctx context.Context) (Session, error)
}
