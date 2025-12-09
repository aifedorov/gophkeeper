package validator

import "fmt"

const (
	LoginMinLength    = 3
	LoginMaxLength    = 25
	PasswordMinLength = 6
	PasswordMaxLength = 16
)

// ValidateLogin validates a login string.
// Login must be between LoginMinLength (3) and LoginMaxLength (25) characters.
func ValidateLogin(login string) error {
	if len(login) == 0 {
		return fmt.Errorf("login can't be empty")
	}
	if len(login) < LoginMinLength {
		return fmt.Errorf("login must be at least %d characters", LoginMinLength)
	}
	if len(login) > LoginMaxLength {
		return fmt.Errorf("login can't be longer than %d characters", LoginMaxLength)
	}
	return nil
}

// ValidatePassword validates a password string.
// Password must be between PasswordMinLength (6) and PasswordMaxLength (16) characters.
func ValidatePassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("password can't be empty")
	}
	if len(password) < PasswordMinLength {
		return fmt.Errorf("password must be at least %d characters", PasswordMinLength)
	}
	if len(password) > PasswordMaxLength {
		return fmt.Errorf("password can't be longer than %d characters", PasswordMaxLength)
	}
	return nil
}
