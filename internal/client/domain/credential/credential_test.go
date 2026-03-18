package credential

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       string
		credName string
		login    string
		password string
		notes    string
		version  int64
	}{
		{
			name:     "creates credential with all fields",
			id:       testID,
			credName: testName,
			login:    testLogin,
			password: testPassword,
			notes:    testNotes,
			version:  testVersion,
		},
		{
			name:     "creates credential without notes",
			id:       "id-2",
			credName: "cred-2",
			login:    "user2",
			password: "pass2",
			notes:    "",
			version:  2,
		},
		{
			name:     "creates credential with zero version",
			id:       "id-3",
			credName: "cred-3",
			login:    "user3",
			password: "pass3",
			notes:    "some notes",
			version:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cred, err := NewCredential(tt.id, tt.credName, tt.login, tt.password, tt.notes, tt.version)

			require.NoError(t, err)
			require.NotNil(t, cred)
			assert.Equal(t, tt.id, cred.ID)
			assert.Equal(t, tt.credName, cred.Name)
			assert.Equal(t, tt.login, cred.Login)
			assert.Equal(t, tt.password, cred.Password)
			assert.Equal(t, tt.notes, cred.Notes)
			assert.Equal(t, tt.version, cred.Version)
		})
	}
}

func TestCredential_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cred    Credential
		wantErr error
	}{
		{
			name: "valid credential with all fields",
			cred: Credential{
				ID:       testID,
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			wantErr: nil,
		},
		{
			name: "valid credential without notes",
			cred: Credential{
				ID:       testID,
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    "",
			},
			wantErr: nil,
		},
		{
			name: "empty id",
			cred: Credential{
				ID:       "",
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
			},
			wantErr: ErrIDRequired,
		},
		{
			name: "empty name",
			cred: Credential{
				ID:       testID,
				Name:     "",
				Login:    testLogin,
				Password: testPassword,
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "empty login",
			cred: Credential{
				ID:       testID,
				Name:     testName,
				Login:    "",
				Password: testPassword,
			},
			wantErr: ErrLoginRequired,
		},
		{
			name: "empty password",
			cred: Credential{
				ID:       testID,
				Name:     testName,
				Login:    testLogin,
				Password: "",
			},
			wantErr: ErrPasswordRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cred.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
