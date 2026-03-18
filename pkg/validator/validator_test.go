package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		login   string
		wantErr bool
	}{
		{
			name:    "valid login",
			login:   "testuser",
			wantErr: false,
		},
		{
			name:    "valid minimum length",
			login:   "abc",
			wantErr: false,
		},
		{
			name:    "valid maximum length",
			login:   "abcdefghij1234567890abcde", // 25 chars
			wantErr: false,
		},
		{
			name:    "login too short",
			login:   "ab",
			wantErr: true,
		},
		{
			name:    "login too long",
			login:   "abcdefghij1234567890abcdef", // 26 chars
			wantErr: true,
		},
		{
			name:    "login empty",
			login:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateLogin(tt.login)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpass",
			wantErr:  false,
		},
		{
			name:     "valid minimum length",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "valid maximum length",
			password: "123456789012345678901234567890", // 30 chars
			wantErr:  false,
		},
		{
			name:     "password too short - less than 6",
			password: "ab",
			wantErr:  true,
		},
		{
			name:     "password too short - exactly 5",
			password: "12345",
			wantErr:  true,
		},
		{
			name:     "password too long",
			password: "1234567890123456789012345678901", // 31 chars
			wantErr:  true,
		},
		{
			name:     "password empty",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidatePassword(tt.password)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSalt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		salt    string
		wantErr bool
	}{
		{
			name:    "valid salt",
			salt:    "somesalt",
			wantErr: false,
		},
		{
			name:    "empty salt",
			salt:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSalt(tt.salt)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
