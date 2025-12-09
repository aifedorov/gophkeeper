package credential

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	id        uuid.UUID
	name      string
	login     string
	password  string
	metadata  string
	createdAt time.Time
	updatedAt time.Time
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
		id:        uuid.New(),
		name:      name,
		login:     login,
		password:  password,
		metadata:  metadata,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}
