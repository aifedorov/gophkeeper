// Package auth provides mappers for converting between domain and repository representations.
package auth

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"github.com/google/uuid"
)

// toDomainUser converts a repository user representation to a domain User entity.
// It parses the user ID from string to UUID and validates the user data.
// Returns an error if the user ID is invalid or user creation fails.
func toDomainUser(user interfaces.RepositoryUser) (User, error) {
	id, err := uuid.Parse(user.ID)
	if err != nil {
		return User{}, fmt.Errorf("failed to parse user id: %w", err)
	}
	res, err := NewUserWithID(id, user.Login, user.Salt)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return res, nil
}

// toRepositoryUser converts a domain User entity to a repository user representation.
// It includes the password hash for storage in the repository.
func toRepositoryUser(user User, passwordHash string) interfaces.RepositoryUser {
	return interfaces.RepositoryUser{
		ID:           user.GetUserID(),
		Login:        user.GetLogin(),
		PasswordHash: passwordHash,
		Salt:         user.GetSalt(),
	}
}
