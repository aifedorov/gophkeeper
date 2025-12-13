package credential

import (
	"testing"

	authMocks "github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces/mocks"
	credMocks "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// Test constants
const (
	testName     = "test-credential"
	testLogin    = "testuser"
	testPassword = "testpassword"
	testNotes    = "test notes"
)

var (
	testUserID = uuid.New()
	testKey    = []byte("test-encryption-key-32-bytes!!")
)

type testSetup struct {
	ctrl           *gomock.Controller
	mockRepo       *credMocks.MockRepository
	mockCrypto     *credMocks.MockCryptoService
	mockSession    *authMocks.MockSessionStore
	service        Service
	logger         *zap.Logger
	userID         uuid.UUID
	credentialID   uuid.UUID
	encryptionKey  []byte
	encryptedLogin []byte
	encryptedPass  []byte
	encryptedNotes []byte
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:           ctrl,
		mockRepo:       credMocks.NewMockRepository(ctrl),
		mockCrypto:     credMocks.NewMockCryptoService(ctrl),
		mockSession:    authMocks.NewMockSessionStore(ctrl),
		logger:         zap.NewNop(),
		userID:         testUserID,
		credentialID:   uuid.New(),
		encryptionKey:  testKey,
		encryptedLogin: []byte("encrypted-login"),
		encryptedPass:  []byte("encrypted-password"),
		encryptedNotes: []byte("encrypted-notes"),
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockRepo, s.mockCrypto, s.mockSession, s.logger)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func (s *testSetup) expectEncryptionKeyInSession() {
	s.mockSession.EXPECT().
		GetEncryptionKey(s.userID).
		Return(s.encryptionKey, true).
		Times(1)
}

func (s *testSetup) expectNoEncryptionKeyInSession() {
	s.mockSession.EXPECT().
		GetEncryptionKey(s.userID).
		Return(nil, false).
		Times(1)
}

func (s *testSetup) expectEncryptCredential() {
	s.mockCrypto.EXPECT().
		Encrypt(testLogin, s.encryptionKey).
		Return(s.encryptedLogin, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Encrypt(testPassword, s.encryptionKey).
		Return(s.encryptedPass, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Encrypt(testNotes, s.encryptionKey).
		Return(s.encryptedNotes, nil).
		Times(1)
}

func (s *testSetup) expectDecryptCredential() {
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedLogin, s.encryptionKey).
		Return(testLogin, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedPass, s.encryptionKey).
		Return(testPassword, nil).
		Times(1)
	s.mockCrypto.EXPECT().
		Decrypt(s.encryptedNotes, s.encryptionKey).
		Return(testNotes, nil).
		Times(1)
}

func newTestCredential() *Credential {
	cred, _ := NewCredential(testName, testLogin, testPassword, testNotes)
	return cred
}

func assertCredentialFields(t *testing.T, cred *Credential, userID uuid.UUID) {
	t.Helper()
	require.NotNil(t, cred)
	assert.Equal(t, testName, cred.GetName())
	assert.Equal(t, testLogin, cred.GetLogin())
	assert.Equal(t, testPassword, cred.GetPassword())
	assert.Equal(t, testNotes, cred.GetMetadata())
	assert.Equal(t, userID, cred.GetUserID())
}

func assertCredentialFieldsWithID(t *testing.T, cred *Credential, userID, credID uuid.UUID) {
	t.Helper()
	assertCredentialFields(t, cred, userID)
	assert.Equal(t, credID, cred.GetID())
}
