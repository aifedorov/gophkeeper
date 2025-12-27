package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAccessToken   = "test-access-token"
	testEncryptionKey = "test-encryption-key"
	testUserID        = "test-user-id"
	testLogin         = "test-login"
)

func TestNewSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		accessToken   string
		encryptionKey string
		userID        string
		login         string
	}{
		{
			name:          "creates session with all fields",
			accessToken:   testAccessToken,
			encryptionKey: testEncryptionKey,
			userID:        testUserID,
			login:         testLogin,
		},
		{
			name:          "creates session with empty fields",
			accessToken:   "",
			encryptionKey: "",
			userID:        "",
			login:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			session := NewSession(tt.accessToken, tt.encryptionKey, tt.userID, tt.login)

			assert.Equal(t, tt.accessToken, session.AccessToken)
			assert.Equal(t, tt.encryptionKey, session.EncryptionKey)
			assert.Equal(t, tt.userID, session.UserID)
			assert.Equal(t, tt.login, session.Login)
		})
	}
}

func TestSession_GetAccessToken(t *testing.T) {
	t.Parallel()

	session := NewSession(testAccessToken, testEncryptionKey, testUserID, testLogin)
	assert.Equal(t, testAccessToken, session.GetAccessToken())
}

func TestSession_GetEncryptionKey(t *testing.T) {
	t.Parallel()

	session := NewSession(testAccessToken, testEncryptionKey, testUserID, testLogin)
	assert.Equal(t, testEncryptionKey, session.GetEncryptionKey())
}

func TestSession_GetUserID(t *testing.T) {
	t.Parallel()

	session := NewSession(testAccessToken, testEncryptionKey, testUserID, testLogin)
	assert.Equal(t, testUserID, session.GetUserID())
}

func TestSession_GetLogin(t *testing.T) {
	t.Parallel()

	session := NewSession(testAccessToken, testEncryptionKey, testUserID, testLogin)
	assert.Equal(t, testLogin, session.GetLogin())
}
