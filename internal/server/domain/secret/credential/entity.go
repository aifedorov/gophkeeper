package credential

import (
	"github.com/google/uuid"
)

type Credential struct {
	id       uuid.UUID
	userID   uuid.UUID
	name     string
	login    string
	password string
	notes    string
}

func NewCredential(name, login, password, metadata string) (*Credential, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if login == "" {
		return nil, ErrLoginRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}

	return &Credential{
		id:       uuid.New(),
		name:     name,
		login:    login,
		password: password,
		notes:    metadata,
	}, nil
}

func (c *Credential) GetID() uuid.UUID {
	return c.id
}

func (c *Credential) GetUserID() uuid.UUID {
	return c.userID
}

func (c *Credential) GetName() string {
	return c.name
}

func (c *Credential) GetLogin() string {
	return c.login
}

func (c *Credential) GetPassword() string {
	return c.password
}

func (c *Credential) GetMetadata() string {
	return c.notes
}
