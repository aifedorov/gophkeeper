package card

import (
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces"
	cardMocks "github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestToDomainCard(t *testing.T) {
	t.Parallel()

	validID := uuid.New()
	validUserID := uuid.New()
	encryptionKey := []byte("test-key-32-bytes-long-string!!")

	tests := []struct {
		name      string
		repoCard  interfaces.RepositoryCard
		setupMock func(*cardMocks.MockCryptoService)
		wantErr   bool
		checkFunc func(*testing.T, Card)
	}{
		{
			name: "successful conversion",
			repoCard: interfaces.RepositoryCard{
				ID:                      validID.String(),
				UserID:                  validUserID.String(),
				Name:                    "test-name",
				EncryptedNumber:         []byte("encrypted-number"),
				EncryptedExpiredDate:    []byte("encrypted-expired-date"),
				EncryptedCardHolderName: []byte("encrypted-card-holder-name"),
				EncryptedCvv:            []byte("encrypted-cvv"),
				EncryptedNotes:          []byte("encrypted-notes"),
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-number"), encryptionKey).
					Return("decrypted-number", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-expired-date"), encryptionKey).
					Return("decrypted-expired-date", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-card-holder-name"), encryptionKey).
					Return("decrypted-card-holder-name", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-cvv"), encryptionKey).
					Return("decrypted-cvv", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), encryptionKey).
					Return("decrypted-notes", nil).
					Times(1)
			},
			checkFunc: func(t *testing.T, c Card) {
				assert.Equal(t, validID.String(), c.GetID())
				assert.Equal(t, validUserID.String(), c.GetUserID())
				assert.Equal(t, "test-name", c.GetName())
				assert.Equal(t, "decrypted-number", c.GetNumber())
				assert.Equal(t, "decrypted-expired-date", c.GetExpiredDate())
				assert.Equal(t, "decrypted-card-holder-name", c.GetCardHolderName())
				assert.Equal(t, "decrypted-cvv", c.GetCvv())
				assert.Equal(t, "decrypted-notes", c.GetNotes())
			},
		},
		{
			name: "number decryption fails",
			repoCard: interfaces.RepositoryCard{
				ID:                      validID.String(),
				UserID:                  validUserID.String(),
				Name:                    "test-name",
				EncryptedNumber:         []byte("encrypted-number"),
				EncryptedExpiredDate:    []byte("encrypted-expired-date"),
				EncryptedCardHolderName: []byte("encrypted-card-holder-name"),
				EncryptedCvv:            []byte("encrypted-cvv"),
				EncryptedNotes:          []byte("encrypted-notes"),
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-number"), encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "expired date decryption fails",
			repoCard: interfaces.RepositoryCard{
				ID:                      validID.String(),
				UserID:                  validUserID.String(),
				Name:                    "test-name",
				EncryptedNumber:         []byte("encrypted-number"),
				EncryptedExpiredDate:    []byte("encrypted-expired-date"),
				EncryptedCardHolderName: []byte("encrypted-card-holder-name"),
				EncryptedCvv:            []byte("encrypted-cvv"),
				EncryptedNotes:          []byte("encrypted-notes"),
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-number"), encryptionKey).
					Return("decrypted-number", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-expired-date"), encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "card holder name decryption fails",
			repoCard: interfaces.RepositoryCard{
				ID:                      validID.String(),
				UserID:                  validUserID.String(),
				Name:                    "test-name",
				EncryptedNumber:         []byte("encrypted-number"),
				EncryptedExpiredDate:    []byte("encrypted-expired-date"),
				EncryptedCardHolderName: []byte("encrypted-card-holder-name"),
				EncryptedCvv:            []byte("encrypted-cvv"),
				EncryptedNotes:          []byte("encrypted-notes"),
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-number"), encryptionKey).
					Return("decrypted-number", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-expired-date"), encryptionKey).
					Return("decrypted-expired-date", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-card-holder-name"), encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "cvv decryption fails",
			repoCard: interfaces.RepositoryCard{
				ID:                      validID.String(),
				UserID:                  validUserID.String(),
				Name:                    "test-name",
				EncryptedNumber:         []byte("encrypted-number"),
				EncryptedExpiredDate:    []byte("encrypted-expired-date"),
				EncryptedCardHolderName: []byte("encrypted-card-holder-name"),
				EncryptedCvv:            []byte("encrypted-cvv"),
				EncryptedNotes:          []byte("encrypted-notes"),
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-number"), encryptionKey).
					Return("decrypted-number", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-expired-date"), encryptionKey).
					Return("decrypted-expired-date", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-card-holder-name"), encryptionKey).
					Return("decrypted-card-holder-name", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-cvv"), encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "notes decryption fails",
			repoCard: interfaces.RepositoryCard{
				ID:                      validID.String(),
				UserID:                  validUserID.String(),
				Name:                    "test-name",
				EncryptedNumber:         []byte("encrypted-number"),
				EncryptedExpiredDate:    []byte("encrypted-expired-date"),
				EncryptedCardHolderName: []byte("encrypted-card-holder-name"),
				EncryptedCvv:            []byte("encrypted-cvv"),
				EncryptedNotes:          []byte("encrypted-notes"),
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-number"), encryptionKey).
					Return("decrypted-number", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-expired-date"), encryptionKey).
					Return("decrypted-expired-date", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-card-holder-name"), encryptionKey).
					Return("decrypted-card-holder-name", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-cvv"), encryptionKey).
					Return("decrypted-cvv", nil).
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

			mockCrypto := cardMocks.NewMockCryptoService(ctrl)
			tt.setupMock(mockCrypto)

			result, err := toDomainCard(mockCrypto, encryptionKey, tt.repoCard)

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

func TestToRepositoryCard(t *testing.T) {
	t.Parallel()

	encryptionKey := []byte("test-key-32-bytes-long-string!!")

	tests := []struct {
		name      string
		card      Card
		setupMock func(*cardMocks.MockCryptoService)
		wantErr   bool
		checkFunc func(*testing.T, interfaces.RepositoryCard)
	}{
		{
			name: "successful conversion",
			card: Card{
				id:             uuid.New().String(),
				userID:         uuid.New().String(),
				name:           "test-name",
				number:         "test-number",
				expiredDate:    "test-expired-date",
				cardHolderName: "test-card-holder-name",
				cvv:            "test-cvv",
				notes:          "test-notes",
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-number", encryptionKey).
					Return([]byte("encrypted-number"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-expired-date", encryptionKey).
					Return([]byte("encrypted-expired-date"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-card-holder-name", encryptionKey).
					Return([]byte("encrypted-card-holder-name"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-cvv", encryptionKey).
					Return([]byte("encrypted-cvv"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-notes", encryptionKey).
					Return([]byte("encrypted-notes"), nil).
					Times(1)
			},
			checkFunc: func(t *testing.T, rc interfaces.RepositoryCard) {
				assert.Equal(t, "test-name", rc.Name)
				assert.Equal(t, []byte("encrypted-number"), rc.EncryptedNumber)
				assert.Equal(t, []byte("encrypted-expired-date"), rc.EncryptedExpiredDate)
				assert.Equal(t, []byte("encrypted-card-holder-name"), rc.EncryptedCardHolderName)
				assert.Equal(t, []byte("encrypted-cvv"), rc.EncryptedCvv)
				assert.Equal(t, []byte("encrypted-notes"), rc.EncryptedNotes)
			},
		},
		{
			name: "number encryption fails",
			card: Card{
				id:             uuid.New().String(),
				userID:         uuid.New().String(),
				name:           "test-name",
				number:         "test-number",
				expiredDate:    "test-expired-date",
				cardHolderName: "test-card-holder-name",
				cvv:            "test-cvv",
				notes:          "test-notes",
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-number", encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "expired date encryption fails",
			card: Card{
				id:             uuid.New().String(),
				userID:         uuid.New().String(),
				name:           "test-name",
				number:         "test-number",
				expiredDate:    "test-expired-date",
				cardHolderName: "test-card-holder-name",
				cvv:            "test-cvv",
				notes:          "test-notes",
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-number", encryptionKey).
					Return([]byte("encrypted-number"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-expired-date", encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "card holder name encryption fails",
			card: Card{
				id:             uuid.New().String(),
				userID:         uuid.New().String(),
				name:           "test-name",
				number:         "test-number",
				expiredDate:    "test-expired-date",
				cardHolderName: "test-card-holder-name",
				cvv:            "test-cvv",
				notes:          "test-notes",
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-number", encryptionKey).
					Return([]byte("encrypted-number"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-expired-date", encryptionKey).
					Return([]byte("encrypted-expired-date"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-card-holder-name", encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "cvv encryption fails",
			card: Card{
				id:             uuid.New().String(),
				userID:         uuid.New().String(),
				name:           "test-name",
				number:         "test-number",
				expiredDate:    "test-expired-date",
				cardHolderName: "test-card-holder-name",
				cvv:            "test-cvv",
				notes:          "test-notes",
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-number", encryptionKey).
					Return([]byte("encrypted-number"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-expired-date", encryptionKey).
					Return([]byte("encrypted-expired-date"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-card-holder-name", encryptionKey).
					Return([]byte("encrypted-card-holder-name"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-cvv", encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "notes encryption fails",
			card: Card{
				id:             uuid.New().String(),
				userID:         uuid.New().String(),
				name:           "test-name",
				number:         "test-number",
				expiredDate:    "test-expired-date",
				cardHolderName: "test-card-holder-name",
				cvv:            "test-cvv",
				notes:          "test-notes",
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-number", encryptionKey).
					Return([]byte("encrypted-number"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-expired-date", encryptionKey).
					Return([]byte("encrypted-expired-date"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-card-holder-name", encryptionKey).
					Return([]byte("encrypted-card-holder-name"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-cvv", encryptionKey).
					Return([]byte("encrypted-cvv"), nil).
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
			card: Card{
				id:             uuid.New().String(),
				userID:         uuid.New().String(),
				name:           "test-name",
				number:         "test-number",
				expiredDate:    "test-expired-date",
				cardHolderName: "test-card-holder-name",
				cvv:            "test-cvv",
				notes:          "",
			},
			setupMock: func(m *cardMocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("test-number", encryptionKey).
					Return([]byte("encrypted-number"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-expired-date", encryptionKey).
					Return([]byte("encrypted-expired-date"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-card-holder-name", encryptionKey).
					Return([]byte("encrypted-card-holder-name"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("test-cvv", encryptionKey).
					Return([]byte("encrypted-cvv"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("", encryptionKey).
					Return([]byte(""), nil).
					Times(1)
			},
			checkFunc: func(t *testing.T, rc interfaces.RepositoryCard) {
				assert.Equal(t, "test-name", rc.Name)
				assert.Equal(t, []byte("encrypted-number"), rc.EncryptedNumber)
				assert.Equal(t, []byte("encrypted-expired-date"), rc.EncryptedExpiredDate)
				assert.Equal(t, []byte("encrypted-card-holder-name"), rc.EncryptedCardHolderName)
				assert.Equal(t, []byte("encrypted-cvv"), rc.EncryptedCvv)
				assert.Equal(t, []byte(""), rc.EncryptedNotes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCrypto := cardMocks.NewMockCryptoService(ctrl)
			tt.setupMock(mockCrypto)

			result, err := toRepositoryCard(mockCrypto, encryptionKey, tt.card)

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
