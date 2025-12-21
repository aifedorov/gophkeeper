// Package binary provides binary file domain errors.
package binary

import "errors"

var (
	// ErrNameExists is returned when attempting to upload a file with a name that already exists for the user.
	ErrNameExists = errors.New("file with this name already exists")
	// ErrNotFound is returned when attempting to access a file that doesn't exist or doesn't belong to the user.
	ErrNotFound = errors.New("file not found")
)
