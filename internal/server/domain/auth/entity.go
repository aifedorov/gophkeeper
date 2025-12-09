package auth

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/pkg/validator"
	"github.com/google/uuid"
)

type User struct {
	id    uuid.UUID
	login string
}

func NewUser(login string) (*User, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return nil, fmt.Errorf("invalid login: %w", err)
	}

	return &User{
		id:    uuid.New(),
		login: login,
	}, nil
}

func (u *User) GetUserID() string {
	return u.id.String()
}

func (u *User) GetLogin() string {
	return u.login
}
