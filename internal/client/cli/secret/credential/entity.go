package credential

import "fmt"

type inputCredentials struct {
	id       string
	name     string
	login    string
	password string
	notes    string
}

func (i *inputCredentials) Validate() error {
	if i.id == "" {
		return fmt.Errorf("id is required")
	}
	if i.name == "" {
		return fmt.Errorf("name is required")
	}
	if i.login == "" {
		return fmt.Errorf("login is required")
	}
	if i.password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
