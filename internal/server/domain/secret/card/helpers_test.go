package card

import (
	"testing"

	cardMocks "github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	testID             = "test-card-id"
	testName           = "test-card"
	testNumber         = "1234567890123456"
	testExpiredDate    = "12/25"
	testCardHolderName = "John Doe"
	testCvv            = "123"
	testNotes          = "test notes"
)

var (
	testUserID        = uuid.New()
	testKey           = []byte("test-encryption-key-32-bytes!!")
	testEncryptionKey = "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEh"
)

type testSetup struct {
	ctrl                    *gomock.Controller
	mockRepo                *cardMocks.MockRepository
	mockCrypto              *cardMocks.MockCryptoService
	service                 Service
	logger                  *zap.Logger
	userID                  string
	cardID                  string
	encryptionKey           []byte
	encryptionKeyStr        string
	encryptedNumber         []byte
	encryptedExpiredDate    []byte
	encryptedCardHolderName []byte
	encryptedCvv            []byte
	encryptedNotes          []byte
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:                    ctrl,
		mockRepo:                cardMocks.NewMockRepository(ctrl),
		mockCrypto:              cardMocks.NewMockCryptoService(ctrl),
		logger:                  zap.NewNop(),
		userID:                  testUserID.String(),
		cardID:                  uuid.New().String(),
		encryptionKey:           testKey,
		encryptionKeyStr:        testEncryptionKey,
		encryptedNumber:         []byte("encrypted-number"),
		encryptedExpiredDate:    []byte("encrypted-expired-date"),
		encryptedCardHolderName: []byte("encrypted-card-holder-name"),
		encryptedCvv:            []byte("encrypted-cvv"),
		encryptedNotes:          []byte("encrypted-notes"),
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockRepo, s.mockCrypto, s.logger)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func (s *testSetup) expectEncryptCard() {
	s.mockCrypto.EXPECT().
		Encrypt(testNumber, s.encryptionKey).
		Return(s.encryptedNumber, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Encrypt(testExpiredDate, s.encryptionKey).
		Return(s.encryptedExpiredDate, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Encrypt(testCardHolderName, s.encryptionKey).
		Return(s.encryptedCardHolderName, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Encrypt(testCvv, s.encryptionKey).
		Return(s.encryptedCvv, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Encrypt(testNotes, s.encryptionKey).
		Return(s.encryptedNotes, nil).
		Times(1)
}

func (s *testSetup) expectDecryptCard() {
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedNumber, s.encryptionKey).
		Return(testNumber, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedExpiredDate, s.encryptionKey).
		Return(testExpiredDate, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedCardHolderName, s.encryptionKey).
		Return(testCardHolderName, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedCvv, s.encryptionKey).
		Return(testCvv, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedNotes, s.encryptionKey).
		Return(testNotes, nil).
		Times(1)
}

func newTestCard() *Card {
	card, _ := NewCard(testID, testName, testNumber, testExpiredDate, testCardHolderName, testCvv, testNotes, 1)
	return card
}

func assertCardFields(t *testing.T, card *Card, userID string) {
	t.Helper()
	require.NotNil(t, card)
	assert.Equal(t, testName, card.GetName())
	assert.Equal(t, testNumber, card.GetNumber())
	assert.Equal(t, testExpiredDate, card.GetExpiredDate())
	assert.Equal(t, testCardHolderName, card.GetCardHolderName())
	assert.Equal(t, testCvv, card.GetCvv())
	assert.Equal(t, testNotes, card.GetNotes())
	assert.Equal(t, userID, card.GetUserID())
}

func assertCardFieldsWithID(t *testing.T, card *Card, userID, cardID string) {
	t.Helper()
	assertCardFields(t, card, userID)
	assert.Equal(t, cardID, card.GetID())
}
