package auth

import (
	"strings"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDomainUser(t *testing.T) {
	t.Parallel()

	validID := uuid.New()
	validSalt := "valid-salt-32-bytes-long-string!!"

	tests := []struct {
		name      string
		repoUser  interfaces.RepositoryUser
		wantErr   bool
		errMsg    string
		checkFunc func(*testing.T, User)
	}{
		{
			name: "valid repository user",
			repoUser: interfaces.RepositoryUser{
				ID:           validID.String(),
				Login:        "testuser",
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, u User) {
				assert.Equal(t, validID.String(), u.GetUserID())
				assert.Equal(t, "testuser", u.GetLogin())
				assert.Equal(t, validSalt, u.GetSalt())
			},
		},
		{
			name: "valid user with different credentials",
			repoUser: interfaces.RepositoryUser{
				ID:           uuid.New().String(),
				Login:        "anotheruser",
				PasswordHash: "another-hash",
				Salt:         "another-salt-32-bytes-long-str!!",
			},
			wantErr: false,
			checkFunc: func(t *testing.T, u User) {
				assert.Equal(t, "anotheruser", u.GetLogin())
				assert.Equal(t, "another-salt-32-bytes-long-str!!", u.GetSalt())
			},
		},
		{
			name: "invalid user ID - not a UUID",
			repoUser: interfaces.RepositoryUser{
				ID:           "invalid-uuid",
				Login:        "testuser",
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: true,
			errMsg:  "failed to parse user id",
		},
		{
			name: "invalid user ID - empty string",
			repoUser: interfaces.RepositoryUser{
				ID:           "",
				Login:        "testuser",
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: true,
			errMsg:  "failed to parse user id",
		},
		{
			name: "invalid user ID - random string",
			repoUser: interfaces.RepositoryUser{
				ID:           "not-a-uuid-at-all",
				Login:        "testuser",
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: true,
			errMsg:  "failed to parse user id",
		},
		{
			name: "invalid login - empty",
			repoUser: interfaces.RepositoryUser{
				ID:           validID.String(),
				Login:        "",
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: true,
			errMsg:  "failed to create user",
		},
		{
			name: "invalid login - too short",
			repoUser: interfaces.RepositoryUser{
				ID:           validID.String(),
				Login:        "ab",
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: true,
			errMsg:  "failed to create user",
		},
		{
			name: "invalid login - too long",
			repoUser: interfaces.RepositoryUser{
				ID:           validID.String(),
				Login:        strings.Repeat("a", 51),
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: true,
			errMsg:  "failed to create user",
		},
		{
			name: "invalid salt - empty",
			repoUser: interfaces.RepositoryUser{
				ID:           validID.String(),
				Login:        "testuser",
				PasswordHash: "hashed-password",
				Salt:         "",
			},
			wantErr: true,
			errMsg:  "failed to create user",
		},
		{
			name: "valid user with minimum login length",
			repoUser: interfaces.RepositoryUser{
				ID:           uuid.New().String(),
				Login:        "abc", // minimum 3 chars
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, u User) {
				assert.Equal(t, "abc", u.GetLogin())
			},
		},
		{
			name: "valid user with maximum login length",
			repoUser: interfaces.RepositoryUser{
				ID:           uuid.New().String(),
				Login:        strings.Repeat("a", 25), // maximum 25 chars
				PasswordHash: "hashed-password",
				Salt:         validSalt,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, u User) {
				assert.Equal(t, strings.Repeat("a", 25), u.GetLogin())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user, err := toDomainUser(tt.repoUser)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, user)
				}
			}
		})
	}
}

func TestToRepositoryUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		login        string
		salt         string
		passwordHash string
		checkFunc    func(*testing.T, User, interfaces.RepositoryUser, string)
	}{
		{
			name:         "standard user conversion",
			login:        "testuser",
			salt:         "valid-salt-32-bytes-long-string!!",
			passwordHash: "hashed-password-123",
			checkFunc: func(t *testing.T, user User, repoUser interfaces.RepositoryUser, hash string) {
				assert.Equal(t, user.GetUserID(), repoUser.ID)
				assert.Equal(t, user.GetLogin(), repoUser.Login)
				assert.Equal(t, hash, repoUser.PasswordHash)
				assert.Equal(t, user.GetSalt(), repoUser.Salt)
			},
		},
		{
			name:         "user with different credentials",
			login:        "anotheruser",
			salt:         "another-salt-32-bytes-long-str!!",
			passwordHash: "different-hash-456",
			checkFunc: func(t *testing.T, user User, repoUser interfaces.RepositoryUser, hash string) {
				assert.Equal(t, "anotheruser", repoUser.Login)
				assert.Equal(t, "another-salt-32-bytes-long-str!!", repoUser.Salt)
				assert.Equal(t, "different-hash-456", repoUser.PasswordHash)
			},
		},
		{
			name:         "user with bcrypt hash",
			login:        "bcryptuser",
			salt:         "bcrypt-salt-32-bytes-long-str!!!",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			checkFunc: func(t *testing.T, user User, repoUser interfaces.RepositoryUser, hash string) {
				assert.Equal(t, user.GetLogin(), repoUser.Login)
				assert.Equal(t, hash, repoUser.PasswordHash)
				assert.True(t, strings.HasPrefix(repoUser.PasswordHash, "$2a$"))
			},
		},
		{
			name:         "user with empty password hash",
			login:        "emptypassuser",
			salt:         "empty-salt-32-bytes-long-str!!!!",
			passwordHash: "",
			checkFunc: func(t *testing.T, user User, repoUser interfaces.RepositoryUser, hash string) {
				assert.Empty(t, repoUser.PasswordHash)
			},
		},
		{
			name:         "user with minimum login length",
			login:        "abc",
			salt:         "min-salt-32-bytes-long-string!!!!",
			passwordHash: "hash-for-min-user",
			checkFunc: func(t *testing.T, user User, repoUser interfaces.RepositoryUser, hash string) {
				assert.Equal(t, "abc", repoUser.Login)
				assert.Len(t, repoUser.Login, 3)
			},
		},
		{
			name:         "user with maximum login length",
			login:        strings.Repeat("a", 25),
			salt:         "max-salt-32-bytes-long-string!!!!",
			passwordHash: "hash-for-max-user",
			checkFunc: func(t *testing.T, user User, repoUser interfaces.RepositoryUser, hash string) {
				assert.Equal(t, strings.Repeat("a", 25), repoUser.Login)
				assert.Len(t, repoUser.Login, 25)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user, err := NewUser(tt.login, tt.salt)
			require.NoError(t, err)

			repoUser := toRepositoryUser(user, tt.passwordHash)

			// Verify UUID format
			_, err = uuid.Parse(repoUser.ID)
			require.NoError(t, err, "ID should be a valid UUID")

			if tt.checkFunc != nil {
				tt.checkFunc(t, user, repoUser, tt.passwordHash)
			}
		})
	}
}

func TestToDomainUser_ToRepositoryUser_RoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("roundtrip conversion preserves data", func(t *testing.T) {
		t.Parallel()

		// Create original repository user
		originalID := uuid.New()
		originalRepoUser := interfaces.RepositoryUser{
			ID:           originalID.String(),
			Login:        "roundtripuser",
			PasswordHash: "original-hash",
			Salt:         "roundtrip-salt-32-bytes-long-str!",
		}

		// Convert to domain user
		domainUser, err := toDomainUser(originalRepoUser)
		require.NoError(t, err)

		// Convert back to repository user
		passwordHash := "pass"
		newRepoUser := toRepositoryUser(domainUser, passwordHash)

		// Verify data integrity (except password hash which is explicitly changed)
		assert.Equal(t, originalRepoUser.ID, newRepoUser.ID)
		assert.Equal(t, originalRepoUser.Login, newRepoUser.Login)
		assert.Equal(t, originalRepoUser.Salt, newRepoUser.Salt)
		assert.Equal(t, passwordHash, newRepoUser.PasswordHash)
		assert.NotEqual(t, originalRepoUser.PasswordHash, newRepoUser.PasswordHash)
	})
}
