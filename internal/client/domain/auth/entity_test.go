package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Parallel()

	t.Run("create auth", func(t *testing.T) {
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
