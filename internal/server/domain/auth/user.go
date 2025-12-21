// Package auth provides authentication domain entities.
package auth

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/pkg/validator"
	"github.com/google/uuid"
)

// User represents a user entity in the authentication domain.
// It contains the user's unique identifier, login, and salt used for encryption key derivation.
type User struct {
	id    uuid.UUID
	login string
	salt  string
}

// NewUser creates a new User entity with a generated UUID.
// It validates the login and salt before creating the user.
// Returns an error if validation fails.
func NewUser(login, salt string) (User, error) {
	return NewUserWithID(uuid.New(), login, salt)
}

// NewUserWithID creates a new User entity with the provided UUID.
// It validates the login and salt before creating the user.
// Returns an error if validation fails.
func NewUserWithID(id uuid.UUID, login, salt string) (User, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return User{}, fmt.Errorf("invalid login: %w", err)
	}

	if err := validator.ValidateSalt(salt); err != nil {
		return User{}, fmt.Errorf("invalid salt: %w", err)
	}

	return User{
		id:    id,
		login: login,
		salt:  salt,
	}, nil
}

// GetUserID returns the user's unique identifier as a string.
func (u *User) GetUserID() string {
	return u.id.String()
}

// GetLogin returns the user's login name.
func (u *User) GetLogin() string {
	return u.login
}

// GetSalt returns the salt used for encryption key derivation.
func (u *User) GetSalt() string {
	return u.salt
}
