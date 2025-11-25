package user

import "github.com/google/uuid"

type User struct {
	userID uuid.UUID
	login  string
}

func NewUser(login string) *User {
	return &User{
		userID: uuid.New(),
		login:  login,
	}
}

func (u *User) GetUserID() string {
	return u.userID.String()
}

func (u *User) GetLogin() string {
	return u.login
}
