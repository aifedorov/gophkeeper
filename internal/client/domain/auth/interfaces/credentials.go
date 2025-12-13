package interfaces

type Credentials struct {
	login    string
	password string
}

func NewCredentials(login, password string) Credentials {
	return Credentials{
		login:    login,
		password: password,
	}
}

func (c Credentials) GetLogin() string {
	return c.login
}

func (c Credentials) GetPassword() string {
	return c.password
}
