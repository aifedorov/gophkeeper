package auth

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/pkg/validator"
	"github.com/google/uuid"
)

type User struct {
	id    uuid.UUID
	login string
	salt  string
}

func NewUser(login, salt string) (User, error) {
	return NewUserWithID(uuid.New(), login, salt)
}

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

func (u *User) GetUserID() string {
	return u.id.String()
}

func (u *User) GetLogin() string {
	return u.login
}

func (u *User) GetSalt() string {
	return u.salt
}
