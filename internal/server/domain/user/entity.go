package user

import "github.com/google/uuid"

type User struct {
	id    uuid.UUID
	login string
}

func NewUser(login string) *User {
	return &User{
		id:    uuid.New(),
		login: login,
	}
}

func (u *User) GetUserID() string {
	return u.id.String()
}

func (u *User) GetLogin() string {
	return u.login
}
