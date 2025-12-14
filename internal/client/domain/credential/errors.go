package credential

import "errors"

var (
	// ErrIDRequired indicates that the id is required.
	ErrIDRequired = errors.New("id can't be empty")
	// ErrNameRequired indicates that the name is required.
	ErrNameRequired = errors.New("name can't be empty")
	// ErrLoginRequired indicates that the login is required.
	ErrLoginRequired = errors.New("login can't be empty")
	// ErrPasswordRequired indicates that the password is required.
	ErrPasswordRequired = errors.New("password can't be empty")
)
