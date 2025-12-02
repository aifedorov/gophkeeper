package user

import "errors"

var (
	ErrLoginExists  = errors.New("login already exists")
	ErrUserNotFound = errors.New("auth.proto not found")
)
