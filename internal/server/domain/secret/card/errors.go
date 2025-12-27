package card

import "errors"

var (
	ErrNameExists             = errors.New("card with this name already exists")
	ErrNotFound               = errors.New("card with this name not found")
	ErrNameRequired           = errors.New("name can't be empty")
	ErrNumberRequired         = errors.New("number can't be empty")
	ErrExpiredDateRequired    = errors.New("expired date can't be empty")
	ErrCardHolderNameRequired = errors.New("card holder name can't be empty")
	ErrCvvRequired            = errors.New("cvv can't be empty")
	ErrIDRequired             = errors.New("id can't be empty")
	ErrVersionConflict        = errors.New("version conflict")
)
