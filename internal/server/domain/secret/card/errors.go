// Package card provides card domain errors.
package card

import "errors"

// Domain errors for card operations.
var (
	// ErrNameExists is returned when a card with the same name already exists.
	ErrNameExists = errors.New("card with this name already exists")
	// ErrNotFound is returned when the card doesn't exist.
	ErrNotFound = errors.New("card with this name not found")
	// ErrNameRequired is returned when card name is empty.
	ErrNameRequired = errors.New("name can't be empty")
	// ErrNumberRequired is returned when card number is empty.
	ErrNumberRequired = errors.New("number can't be empty")
	// ErrExpiredDateRequired is returned when expiration date is empty.
	ErrExpiredDateRequired = errors.New("expired date can't be empty")
	// ErrCardHolderNameRequired is returned when cardholder name is empty.
	ErrCardHolderNameRequired = errors.New("card holder name can't be empty")
	// ErrCvvRequired is returned when CVV is empty.
	ErrCvvRequired = errors.New("cvv can't be empty")
	// ErrIDRequired is returned when card ID is empty.
	ErrIDRequired = errors.New("id can't be empty")
	// ErrVersionConflict is returned when card version doesn't match.
	ErrVersionConflict = errors.New("version conflict")
)
