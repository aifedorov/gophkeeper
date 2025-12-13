package credential

import (
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	credMocks "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestToDomainCredential(t *testing.T) {
	t.Parallel()

	validID := uuid.New()
	validUserID := uuid.New()
	encryptionKey := []byte("test-key-32-bytes-long-string!!")

	tests := []struct {
		name      string
		repoCred  interfaces.RepositoryCredential
		setupMock func(*credMocks.MockCryptoService)
		wantErr   bool
		checkFunc func(*testing.T, Credential)
	}{
		{
			name: "successful conversion",
			repoCred: interfaces.RepositoryCredential{
				ID:                validID.String(),
				UserID:            validUserID.String(),
				Name:              "test-name",
				EncryptedLogin:    []byte("encrypted-login"),
				EncryptedPassword: []byte("encrypted-password"),
				EncryptedNotes:    []byte("encrypted-notes"),
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-login"), encryptionKey).
					Return("decrypted-login", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-password"), encryptionKey).
					Return("decrypted-password", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), encryptionKey).
					Return("decrypted-notes", nil).
					Times(1)
			},
			checkFunc: func(t *testing.T, c Credential) {
				assert.Equal(t, validID.String(), c.GetID())
				assert.Equal(t, validUserID.String(), c.GetUserID())
				assert.Equal(t, "test-name", c.GetName())
				assert.Equal(t, "decrypted-login", c.GetLogin())
				assert.Equal(t, "decrypted-password", c.GetPassword())
				assert.Equal(t, "decrypted-notes", c.GetMetadata())
			},
		},
		{
			name: "login decryption fails",
			repoCred: interfaces.RepositoryCredential{
				ID:                validID.String(),
				UserID:            validUserID.String(),
				Name:              "test-name",
				EncryptedLogin:    []byte("encrypted-login"),
				EncryptedPassword: []byte("encrypted-password"),
				EncryptedNotes:    []byte("encrypted-notes"),
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-login"), encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "password decryption fails",
			repoCred: interfaces.RepositoryCredential{
				ID:                validID.String(),
				UserID:            validUserID.String(),
				Name:              "test-name",
				EncryptedLogin:    []byte("encrypted-login"),
				EncryptedPassword: []byte("encrypted-password"),
				EncryptedNotes:    []byte("encrypted-notes"),
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-login"), encryptionKey).
					Return("decrypted-login", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-password"), encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "notes decryption fails",
			repoCred: interfaces.RepositoryCredential{
				ID:                validID.String(),
				UserID:            validUserID.String(),
				Name:              "test-name",
				EncryptedLogin:    []byte("encrypted-login"),
				EncryptedPassword: []byte("encrypted-password"),
				EncryptedNotes:    []byte("encrypted-notes"),
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-login"), encryptionKey).
					Return("decrypted-login", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-password"), encryptionKey).
					Return("decrypted-password", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCrypto := credMocks.NewMockCryptoService(ctrl)
			tt.setupMock(mockCrypto)

			result, err := toDomainCredential(mockCrypto, encryptionKey, tt.repoCred)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}

func TestToRepositoryCredential(t *testing.T) {
	t.Parallel()

	encryptionKey := []byte("test-key-32-bytes-long-string!!")

	tests := []struct {
		name      string
		cred      Credential
		setupMock func(*credMocks.MockCryptoService)
		wantErr   bool
		checkFunc func(*testing.T, interfaces.RepositoryCredential)
	}{
		{
			name: "successful conversion",
			cred: Credential{
				id:       uuid.New().String(),
				userID:   uuid.New().String(),
				name:     "test-name",
				login:    "test-login",
				password: "test-password",
				notes:    "test-notes",
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-login", encryptionKey).
					Return([]byte("encrypted-login"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-password", encryptionKey).
					Return([]byte("encrypted-password"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-notes", encryptionKey).
					Return([]byte("encrypted-notes"), nil).
					Times(1)
			},
			checkFunc: func(t *testing.T, rc interfaces.RepositoryCredential) {
				assert.Equal(t, "test-name", rc.Name)
				assert.Equal(t, []byte("encrypted-login"), rc.EncryptedLogin)
				assert.Equal(t, []byte("encrypted-password"), rc.EncryptedPassword)
				assert.Equal(t, []byte("encrypted-notes"), rc.EncryptedNotes)
			},
		},
		{
			name: "login encryption fails",
			cred: Credential{
				id:       uuid.New().String(),
				userID:   uuid.New().String(),
				name:     "test-name",
				login:    "test-login",
				password: "test-password",
				notes:    "test-notes",
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-login", encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "password encryption fails",
			cred: Credential{
				id:       uuid.New().String(),
				userID:   uuid.New().String(),
				name:     "test-name",
				login:    "test-login",
				password: "test-password",
				notes:    "test-notes",
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-login", encryptionKey).
					Return([]byte("encrypted-login"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-password", encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "notes encryption fails",
			cred: Credential{
				id:       uuid.New().String(),
				userID:   uuid.New().String(),
				name:     "test-name",
				login:    "test-login",
				password: "test-password",
				notes:    "test-notes",
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-login", encryptionKey).
					Return([]byte("encrypted-login"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-password", encryptionKey).
					Return([]byte("encrypted-password"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-notes", encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "empty notes",
			cred: Credential{
				id:       uuid.New().String(),
				userID:   uuid.New().String(),
				name:     "test-name",
				login:    "test-login",
				password: "test-password",
				notes:    "",
			},
			setupMock: func(m *credMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-login", encryptionKey).
					Return([]byte("encrypted-login"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-password", encryptionKey).
					Return([]byte("encrypted-password"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("", encryptionKey).
					Return([]byte(""), nil).
					Times(1)
			},
			checkFunc: func(t *testing.T, rc interfaces.RepositoryCredential) {
				assert.Equal(t, "test-name", rc.Name)
				assert.Equal(t, []byte("encrypted-login"), rc.EncryptedLogin)
				assert.Equal(t, []byte("encrypted-password"), rc.EncryptedPassword)
				assert.Equal(t, []byte(""), rc.EncryptedNotes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCrypto := credMocks.NewMockCryptoService(ctrl)
			tt.setupMock(mockCrypto)

			result, err := toRepositoryCredential(mockCrypto, encryptionKey, tt.cred)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}
