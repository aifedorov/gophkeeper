package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		login    string
		password string
	}{
		{
			name:     "creates credentials with all fields",
			login:    "testuser",
			password: "testpass",
		},
		{
			name:     "creates credentials with empty fields",
			login:    "",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			creds := NewCredentials(tt.login, tt.password)

			assert.Equal(t, tt.login, creds.GetLogin())
			assert.Equal(t, tt.password, creds.GetPassword())
		})
	}
}
