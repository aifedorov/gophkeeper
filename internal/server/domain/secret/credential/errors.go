package credential

import "errors"

var (
	ErrNameExists       = errors.New("credential with this name already exists")
	ErrNotFound         = errors.New("credential with this name not found")
	ErrNameRequired     = errors.New("name can't be empty")
	ErrLoginRequired    = errors.New("login can't be empty")
	ErrPasswordRequired = errors.New("password can't be empty")
)
