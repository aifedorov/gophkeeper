// Package interfaces provides authentication credentials for the GophKeeper client.
package interfaces

// Credentials represents user login credentials.
// It contains the login (username/email) and password.
type Credentials struct {
	login    string
	password string
}

// NewCredentials creates a new Credentials instance with the provided login and password.
func NewCredentials(login, password string) Credentials {
	return Credentials{
		login:    login,
		password: password,
	}
}

// GetLogin returns the login (username or email).
func (c Credentials) GetLogin() string {
	return c.login
}

// GetPassword returns the password.
func (c Credentials) GetPassword() string {
	return c.password
}
