package validator

import "fmt"

const (
	LoginMinLength    = 3
	LoginMaxLength    = 25
	PasswordMinLength = 6
	PasswordMaxLength = 16
)

func ValidateLogin(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("login can't be empty")
	}
	if len(s) < LoginMinLength {
		return fmt.Errorf("login must be at least %d characters", LoginMinLength)
	}
	if len(s) > LoginMaxLength {
		return fmt.Errorf("login can't be longer than %d characters", LoginMaxLength)
	}
	return nil
}

func ValidatePassword(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("password can't be empty")
	}
	if len(s) < PasswordMinLength {
		return fmt.Errorf("password must be at least %d characters", PasswordMinLength)
	}
	if len(s) > PasswordMaxLength {
		return fmt.Errorf("password can't be longer than %d characters", PasswordMaxLength)
	}
	return nil
}
