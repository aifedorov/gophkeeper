// Package auth provides authentication domain errors.
package auth

import "errors"

var (
	// ErrLoginExists is returned when attempting to register a user with a login that already exists.
	ErrLoginExists = errors.New("login already exists")
	// ErrUserNotFound is returned when attempting to login with credentials for a user that doesn't exist.
	ErrUserNotFound = errors.New("user not found")
)
