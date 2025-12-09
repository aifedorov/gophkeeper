package repository

import "errors"

var (
	ErrNameExists = errors.New("credential with this name already exists")
	ErrNotFound   = errors.New("credential not found")
)
