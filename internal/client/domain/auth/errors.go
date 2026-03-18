package auth

import "errors"

var (
	// ErrInvalidCredentials indicates invalid login or password during authentication.
	ErrInvalidCredentials = errors.New("wrong password or login")

	// ErrUserAlreadyExists indicates that a auth with the specified login already exists.
	ErrUserAlreadyExists = errors.New("auth with this login already exists")

	// ErrInvalidLogin indicates that the login must be between 3 and 25 characters.
	ErrInvalidLogin = errors.New("login mus be from 3 to 25 symbols")

	// ErrInvalidPassword indicates that the password must be between 3 and 16 characters.
	ErrInvalidPassword = errors.New("password must be from 3 to 16 characters")

	// ErrSessionNotFound indicates that the session is not found.
	ErrSessionNotFound = errors.New("session not found")
)
