package repository

import "errors"

var (
	ErrLoginExists  = errors.New("login already exists")
	ErrUserNotFound = errors.New("user not found")
)
