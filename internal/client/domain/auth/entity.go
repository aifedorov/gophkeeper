package auth

type User struct {
	ID    string
	Login string
}

type Session struct {
	User        User
	AccessToken string
}

type Credentials struct {
	Login    string
	Password string
}

func (c *Credentials) Validate() error {
	if len(c.Login) < 3 || len(c.Login) > 25 {
		return ErrInvalidLogin
	}
	if len(c.Password) < 3 || len(c.Password) > 16 {
		return ErrInvalidPassword
	}
	return nil
}
