package card

import "errors"

var (
	// ErrIDRequired indicates that the id is required.
	ErrIDRequired = errors.New("id can't be empty")
	// ErrNameRequired indicates that the name is required.
	ErrNameRequired = errors.New("name can't be empty")
	// ErrNumberRequired indicates that the number is required.
	ErrNumberRequired = errors.New("number can't be empty")
	// ErrExpiredDateRequired indicates that the expired date is required.
	ErrExpiredDateRequired = errors.New("expired date can't be empty")
	// ErrCardHolderNameRequired indicates that the card holder name is required.
	ErrCardHolderNameRequired = errors.New("card holder name can't be empty")
	// ErrCvvRequired indicates that the cvv is required.
	ErrCvvRequired = errors.New("cvv can't be empty")
)
