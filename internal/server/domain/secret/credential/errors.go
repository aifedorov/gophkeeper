// Package credential provides credential domain errors.
package credential

import "errors"

var (
	// ErrNameExists is returned when attempting to create a credential with a name that already exists for the user.
	ErrNameExists = errors.New("credential with this name already exists")
	// ErrNotFound is returned when attempting to access a credential that doesn't exist or doesn't belong to the user.
	ErrNotFound = errors.New("credential with this name not found")
	// ErrNameRequired is returned when attempting to create a credential without a name.
	ErrNameRequired = errors.New("name can't be empty")
	// ErrLoginRequired is returned when attempting to create a credential without a login.
	ErrLoginRequired = errors.New("login can't be empty")
	// ErrPasswordRequired is returned when attempting to create a credential without a password.
	ErrPasswordRequired = errors.New("password can't be empty")
	// ErrIDRequired is returned when attempting to create a credential without an ID.
	ErrIDRequired = errors.New("id can't be empty")
)
