// Package binary provides binary file domain errors.
package binary

import "errors"

var (
	// ErrFileNotFound is returned when attempting to access a binary file that doesn't exist
	// or doesn't belong to the user.
	ErrFileNotFound = errors.New("binary with this name not found")
)
