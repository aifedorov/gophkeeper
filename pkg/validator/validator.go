// Package validator provides validation functions for user input.
package validator

import "fmt"

// Validation constraints for login and password fields.
const (
	// LoginMinLength is the minimum required length for login.
	LoginMinLength = 3
	// LoginMaxLength is the maximum allowed length for login.
	LoginMaxLength = 25
	// PasswordMinLength is the minimum required length for password.
	PasswordMinLength = 6
	// PasswordMaxLength is the maximum allowed length for password.
	PasswordMaxLength = 16
)

// ValidateLogin validates a login string.
// login must be between LoginMinLength (3) and LoginMaxLength (25) characters.
func ValidateLogin(login string) error {
	if len(login) == 0 {
		return fmt.Errorf("validator: login can't be empty")
	}
	if len(login) < LoginMinLength {
		return fmt.Errorf("validator: login must be at least %d characters", LoginMinLength)
	}
	if len(login) > LoginMaxLength {
		return fmt.Errorf("validator: login can't be longer than %d characters", LoginMaxLength)
	}
	return nil
}

// ValidatePassword validates a password string.
// password must be between PasswordMinLength (6) and PasswordMaxLength (16) characters.
func ValidatePassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("validator: password can't be empty")
	}
	if len(password) < PasswordMinLength {
		return fmt.Errorf("validator: password must be at least %d characters", PasswordMinLength)
	}
	if len(password) > PasswordMaxLength {
		return fmt.Errorf("validator: password can't be longer than %d characters", PasswordMaxLength)
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
