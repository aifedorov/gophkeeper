package user

import repository "github.com/aifedorov/gophkeeper/internal/domain/user/repository/db"

func toDomainUser(user *repository.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		user.ID,
		user.Login,
	}
}
