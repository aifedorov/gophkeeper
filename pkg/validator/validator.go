// Package validator provides validation functions for user input.
package validator

import "fmt"

// Validation constraints for login and password fields.
const (
	loginMinLength    = 3
	loginMaxLength    = 25
	passwordMinLength = 6
	passwordMaxLength = 30
)

// ValidateLogin validates a login string.
func ValidateLogin(login string) error {
	if len(login) == 0 {
		return fmt.Errorf("validator: login can't be empty")
	}
	if len(login) < loginMinLength {
		return fmt.Errorf("validator: login must be at least %d characters", loginMinLength)
	}
	if len(login) > loginMaxLength {
		return fmt.Errorf("validator: login can't be longer than %d characters", loginMaxLength)
	}
	return nil
}

// ValidatePassword validates a password string.
func ValidatePassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("validator: password can't be empty")
	}
	if len(password) < passwordMinLength {
		return fmt.Errorf("validator: password must be at least %d characters", passwordMinLength)
	}
	if len(password) > passwordMaxLength {
		return fmt.Errorf("validator: password can't be longer than %d characters", passwordMaxLength)
	}
	return nil
}

// ValidateSalt validates a salt string. Salt must not be empty.
func ValidateSalt(salt string) error {
	if len(salt) == 0 {
		return fmt.Errorf("validator: salt can't be empty")
	}
	return nil
}
