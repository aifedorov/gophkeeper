package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentials_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		credentials Credentials
		wantErr     error
	}{
		{
			name: "valid credentials",
			credentials: Credentials{
				Login:    "testuser",
				Password: "testpass",
			},
			wantErr: nil,
		},
		{
			name: "valid minimum length",
			credentials: Credentials{
				Login:    "abc",
				Password: "123",
			},
			wantErr: nil,
		},
		{
			name: "valid maximum length",
			credentials: Credentials{
				Login:    "abcdefghij1234567890abcde", // 25 chars
				Password: "1234567890123456",          // 16 chars
			},
			wantErr: nil,
		},
		{
			name: "login too short",
			credentials: Credentials{
				Login:    "ab",
				Password: "testpass",
			},
			wantErr: ErrInvalidLogin,
		},
		{
			name: "login too long",
			credentials: Credentials{
				Login:    "abcdefghij1234567890abcdef", // 26 chars
				Password: "testpass",
			},
			wantErr: ErrInvalidLogin,
		},
		{
			name: "login empty",
			credentials: Credentials{
				Login:    "",
				Password: "testpass",
			},
			wantErr: ErrInvalidLogin,
		},
		{
			name: "password too short",
			credentials: Credentials{
				Login:    "testuser",
				Password: "ab",
			},
			wantErr: ErrInvalidPassword,
		},
		{
			name: "password too long",
			credentials: Credentials{
				Login:    "testuser",
				Password: "12345678901234567", // 17 chars
			},
			wantErr: ErrInvalidPassword,
		},
		{
			name: "password empty",
			credentials: Credentials{
				Login:    "testuser",
				Password: "",
			},
			wantErr: ErrInvalidPassword,
		},
		{
			name: "both empty",
			credentials: Credentials{
				Login:    "",
				Password: "",
			},
			wantErr: ErrInvalidLogin, // login checked first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.credentials.Validate()

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUser(t *testing.T) {
	t.Parallel()

	t.Run("create user", func(t *testing.T) {
		t.Parallel()

		user := User{
			ID:    "test-id-123",
			Login: "testuser",
		}

		assert.Equal(t, "test-id-123", user.ID)
		assert.Equal(t, "testuser", user.Login)
	})
}

func TestSession(t *testing.T) {
	t.Parallel()

	t.Run("create session", func(t *testing.T) {
		t.Parallel()

		session := Session{
			User: User{
				ID:    "test-id-123",
				Login: "testuser",
			},
			AccessToken: "test-token-xyz",
		}

		assert.Equal(t, "test-id-123", session.User.ID)
		assert.Equal(t, "testuser", session.User.Login)
		assert.Equal(t, "test-token-xyz", session.AccessToken)
	})
}
