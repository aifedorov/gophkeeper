package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestNewService(t *testing.T) {
	logger := zap.NewNop()
	service := NewService(logger)

	require.NotNil(t, service)
	assert.NotNil(t, service.logger)
}

func TestService_GenerateSalt(t *testing.T) {
	t.Parallel()

	service := NewService(zap.NewNop())

	t.Run("successful salt generation", func(t *testing.T) {
		salt, err := service.GenerateSalt()

		require.NoError(t, err)
		assert.Len(t, salt, saltLen)
	})

	t.Run("generates unique salts", func(t *testing.T) {
		salt1, err1 := service.GenerateSalt()
		salt2, err2 := service.GenerateSalt()

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, salt1, salt2)
	})
}

func TestService_DeriveEncryptionKey(t *testing.T) {
	t.Parallel()

	service := NewService(zap.NewNop())

	tests := []struct {
		name     string
		password string
		salt     string
	}{
		{
			name:     "standard password and salt",
			password: "test-password",
			salt:     "test-salt-32-bytes-long-string!!",
		},
		{
			name:     "empty password",
			password: "",
			salt:     "test-salt-32-bytes-long-string!!",
		},
		{
			name:     "empty salt",
			password: "test-password",
			salt:     "",
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 1000),
			salt:     "test-salt-32-bytes-long-string!!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			key := service.DeriveEncryptionKey(tt.password, tt.salt)

			assert.Len(t, key, argonKeyLen)
		})
	}

	t.Run("same password and salt produce same key", func(t *testing.T) {
		password := "test-password"
		salt := "test-salt-32-bytes-long-string!!"

		key1 := service.DeriveEncryptionKey(password, salt)
		key2 := service.DeriveEncryptionKey(password, salt)

		assert.Equal(t, key1, key2)
	})

	t.Run("different passwords produce different keys", func(t *testing.T) {
		salt := "test-salt-32-bytes-long-string!!"

		key1 := service.DeriveEncryptionKey("password1", salt)
		key2 := service.DeriveEncryptionKey("password2", salt)

		assert.NotEqual(t, key1, key2)
	})

	t.Run("different salts produce different keys", func(t *testing.T) {
		password := "test-password"

		key1 := service.DeriveEncryptionKey(password, "salt1-32-bytes-long-string!!!!!")
		key2 := service.DeriveEncryptionKey(password, "salt2-32-bytes-long-string!!!!!")

		assert.NotEqual(t, key1, key2)
	})
}

func TestService_HashPassword(t *testing.T) {
	t.Parallel()

	service := NewService(zap.NewNop())

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "standard password",
			password: "test-password",
			wantErr:  false,
		},
		{
			name:     "short password",
			password: "abc",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 72),
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "p@ssw0rd!#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "unicode password",
			password: "пароль",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash, err := service.HashPassword(tt.password)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.True(t, strings.HasPrefix(hash, "$2a$"))
			}
		})
	}

	t.Run("same password produces different hashes", func(t *testing.T) {
		password := "test-password"

		hash1, err1 := service.HashPassword(password)
		hash2, err2 := service.HashPassword(password)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2)
	})
}

func TestService_CompareHashAndPassword(t *testing.T) {
	t.Parallel()

	service := NewService(zap.NewNop())

	t.Run("correct password matches hash", func(t *testing.T) {
		password := "test-password"
		hash, err := service.HashPassword(password)
		require.NoError(t, err)

		err = service.CompareHashAndPassword(hash, password)
		assert.NoError(t, err)
	})

	t.Run("incorrect password does not match hash", func(t *testing.T) {
		password := "test-password"
		hash, err := service.HashPassword(password)
		require.NoError(t, err)

		err = service.CompareHashAndPassword(hash, "wrong-password")
		assert.Error(t, err)
		assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
	})

	t.Run("invalid hash format", func(t *testing.T) {
		err := service.CompareHashAndPassword("invalid-hash", "password")
		assert.Error(t, err)
	})

	t.Run("empty password", func(t *testing.T) {
		password := ""
		hash, err := service.HashPassword(password)
		require.NoError(t, err)

		err = service.CompareHashAndPassword(hash, password)
		assert.NoError(t, err)
	})
}

func TestService_Encrypt(t *testing.T) {
	t.Parallel()

	service := NewService(zap.NewNop())
	validKey := make([]byte, aes256KeySize)

	tests := []struct {
		name      string
		plaintext string
		key       []byte
		wantErr   bool
	}{
		{
			name:      "standard text",
			plaintext: "Hello, World!",
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "empty text",
			plaintext: "",
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "long text",
			plaintext: strings.Repeat("a", 10000),
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "unicode text",
			plaintext: "Привет, 世界!",
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "invalid key size - too short",
			plaintext: "test",
			key:       make([]byte, 16),
			wantErr:   true,
		},
		{
			name:      "invalid key size - too long",
			plaintext: "test",
			key:       make([]byte, 64),
			wantErr:   true,
		},
		{
			name:      "empty key",
			plaintext: "test",
			key:       []byte{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ciphertext, err := service.Encrypt(tt.plaintext, tt.key)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, ciphertext)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, ciphertext)
				assert.NotEqual(t, []byte(tt.plaintext), ciphertext)
			}
		})
	}

	t.Run("same plaintext produces different ciphertexts", func(t *testing.T) {
		plaintext := "test-data"

		ciphertext1, err1 := service.Encrypt(plaintext, validKey)
		ciphertext2, err2 := service.Encrypt(plaintext, validKey)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, ciphertext1, ciphertext2)
	})
}

func TestService_Decrypt(t *testing.T) {
	t.Parallel()

	service := NewService(zap.NewNop())
	validKey := make([]byte, aes256KeySize)

	t.Run("decrypt successfully encrypted data", func(t *testing.T) {
		plaintext := "Hello, World!"

		ciphertext, err := service.Encrypt(plaintext, validKey)
		require.NoError(t, err)

		decrypted, err := service.Decrypt(ciphertext, validKey)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("decrypt empty string", func(t *testing.T) {
		plaintext := ""

		ciphertext, err := service.Encrypt(plaintext, validKey)
		require.NoError(t, err)

		decrypted, err := service.Decrypt(ciphertext, validKey)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("decrypt long text", func(t *testing.T) {
		plaintext := strings.Repeat("test data ", 1000)

		ciphertext, err := service.Encrypt(plaintext, validKey)
		require.NoError(t, err)

		decrypted, err := service.Decrypt(ciphertext, validKey)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("decrypt unicode text", func(t *testing.T) {
		plaintext := "Привет, 世界! 🌍"

		ciphertext, err := service.Encrypt(plaintext, validKey)
		require.NoError(t, err)

		decrypted, err := service.Decrypt(ciphertext, validKey)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("wrong key fails decryption", func(t *testing.T) {
		plaintext := "secret data"
		wrongKey := make([]byte, aes256KeySize)
		for i := range wrongKey {
			wrongKey[i] = 0xFF
		}

		ciphertext, err := service.Encrypt(plaintext, validKey)
		require.NoError(t, err)

		_, err = service.Decrypt(ciphertext, wrongKey)
		assert.Error(t, err)
	})

	t.Run("invalid key size", func(t *testing.T) {
		ciphertext := []byte("some-data")
		invalidKey := make([]byte, 16)

		_, err := service.Decrypt(ciphertext, invalidKey)
		assert.Error(t, err)
	})

	t.Run("ciphertext too short", func(t *testing.T) {
		shortCiphertext := []byte("short")

		_, err := service.Decrypt(shortCiphertext, validKey)
		assert.Error(t, err)
	})

	t.Run("corrupted ciphertext", func(t *testing.T) {
		plaintext := "test data"

		ciphertext, err := service.Encrypt(plaintext, validKey)
		require.NoError(t, err)

		// Corrupt the ciphertext
		ciphertext[len(ciphertext)-1] ^= 0xFF

		_, err = service.Decrypt(ciphertext, validKey)
		assert.Error(t, err)
	})

	t.Run("empty ciphertext", func(t *testing.T) {
		_, err := service.Decrypt([]byte{}, validKey)
		assert.Error(t, err)
	})
}

func TestService_EncryptDecrypt_Integration(t *testing.T) {
	t.Parallel()

	service := NewService(zap.NewNop())

	testCases := []string{
		"simple text",
		"",
		"text with spaces and numbers 123",
		"special chars !@#$%^&*()",
		"unicode: Привет мир 世界",
		strings.Repeat("long text ", 100),
		"line1\nline2\nline3",
		"tab\tseparated\tvalues",
	}

	for _, plaintext := range testCases {
		t.Run("roundtrip: "+plaintext[:min(20, len(plaintext))], func(t *testing.T) {
			t.Parallel()

			key := make([]byte, aes256KeySize)

			ciphertext, err := service.Encrypt(plaintext, key)
			require.NoError(t, err)

			decrypted, err := service.Decrypt(ciphertext, key)
			require.NoError(t, err)

			assert.Equal(t, plaintext, decrypted)
		})
	}
}
