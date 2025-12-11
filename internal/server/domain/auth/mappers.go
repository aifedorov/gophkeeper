package auth

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"github.com/google/uuid"
)

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

func toRepositoryUser(user User, passwordHash string) interfaces.RepositoryUser {
	return interfaces.RepositoryUser{
		ID:           user.GetUserID(),
		Login:        user.GetLogin(),
		PasswordHash: passwordHash,
		Salt:         user.GetSalt(),
	}
}
