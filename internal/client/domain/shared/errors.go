package shared

import "errors"

var (
	// ErrVersionConflict indicates that the credential was modified by another client.
	ErrVersionConflict = errors.New("credential was modified by another client, please refetch and retry")
	// ErrNotFound indicates that the credential was not found on the server.
	ErrNotFound = errors.New("credential not found")
	// ErrAlreadyExists indicates that a credential with this name already exists.
	ErrAlreadyExists = errors.New("credential with this name already exists")
	// ErrUnauthenticated indicates that the user is not authenticated or token is invalid.
	ErrUnauthenticated = errors.New("authentication required or token expired")
)
