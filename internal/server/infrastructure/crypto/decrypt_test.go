package crypto

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func encryptData(plaintext string, key []byte) (string, error) {
	plainReader := strings.NewReader(plaintext)
	encryptReader, err := NewEncryptReader(plainReader, key)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	_, err = io.Copy(&buf, encryptReader)
	return buf.String(), err
}

func decryptData(ciphertext string, key []byte) (string, error) {
	cipherReader := strings.NewReader(ciphertext)
	decryptReader, err := NewDecryptReader(cipherReader, key)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	_, err = io.Copy(&buf, decryptReader)
	return buf.String(), err
}

func TestNewDecryptReader(t *testing.T) {
	t.Parallel()

	validKey := make([]byte, aes256KeySize)

	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "valid key - AES-256",
			key:     validKey,
			wantErr: false,
		},
		{
			name:    "valid key - AES-128",
			key:     make([]byte, 16),
			wantErr: false,
		},
		{
			name:    "invalid key - too short",
			key:     make([]byte, 8),
			wantErr: true,
		},
		{
			name:    "empty key",
			key:     []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := strings.NewReader("test data")
			decryptReader, err := NewDecryptReader(reader, tt.key)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, decryptReader)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, decryptReader)
			}
		})
	}
}

func TestDecryptReader_Read(t *testing.T) {
	t.Parallel()

	validKey := make([]byte, aes256KeySize)

	testCases := []struct {
		name string
		data string
	}{
		{"small data", "Hello, World!"},
		{"empty data", ""},
		{"multiple chunks", strings.Repeat("test data ", 10000)},
		{"exactly one chunk", strings.Repeat("x", chunkSize)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encrypted, err := encryptData(tc.data, validKey)
			require.NoError(t, err)

			decrypted, err := decryptData(encrypted, validKey)
			require.NoError(t, err)

			assert.Equal(t, tc.data, decrypted)
		})
	}

	t.Run("missing nonce", func(t *testing.T) {
		reader := strings.NewReader("short")
		decryptReader, err := NewDecryptReader(reader, validKey)
		require.NoError(t, err)

		buf := make([]byte, 100)
		_, err = decryptReader.Read(buf)
		assert.Error(t, err)
	})

	t.Run("corrupted ciphertext", func(t *testing.T) {
		encrypted, err := encryptData("test data", validKey)
		require.NoError(t, err)

		// Corrupt the last byte
		corrupted := []byte(encrypted)
		corrupted[len(corrupted)-1] ^= 0xFF

		decryptReader, err := NewDecryptReader(strings.NewReader(string(corrupted)), validKey)
		require.NoError(t, err)

		buf := make([]byte, 100)
		_, err = decryptReader.Read(buf)
		assert.Error(t, err)
	})

	t.Run("wrong key fails decryption", func(t *testing.T) {
		wrongKey := make([]byte, aes256KeySize)
		for i := range wrongKey {
			wrongKey[i] = 0xFF
		}

		encrypted, err := encryptData("secret data", validKey)
		require.NoError(t, err)

		_, err = decryptData(encrypted, wrongKey)
		assert.Error(t, err)
	})
}
