package credential

import (
	"testing"

	credMocks "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// Test constants
const (
	testID       = "test-credential-id"
	testName     = "test-credential"
	testLogin    = "testuser"
	testPassword = "testpassword"
	testNotes    = "test notes"
)

var (
	testUserID        = uuid.New()
	testKey           = []byte("test-encryption-key-32-bytes!!")
	testEncryptionKey = "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEh" // base64 encoded testKey
)

type testSetup struct {
	ctrl             *gomock.Controller
	mockRepo         *credMocks.MockRepository
	mockCrypto       *credMocks.MockCryptoService
	service          Service
	logger           *zap.Logger
	userID           string
	credentialID     string
	encryptionKey    []byte
	encryptionKeyStr string
	encryptedLogin   []byte
	encryptedPass    []byte
	encryptedNotes   []byte
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:             ctrl,
		mockRepo:         credMocks.NewMockRepository(ctrl),
		mockCrypto:       credMocks.NewMockCryptoService(ctrl),
		logger:           zap.NewNop(),
		userID:           testUserID.String(),
		credentialID:     uuid.New().String(),
		encryptionKey:    testKey,
		encryptionKeyStr: testEncryptionKey,
		encryptedLogin:   []byte("encrypted-login"),
		encryptedPass:    []byte("encrypted-password"),
		encryptedNotes:   []byte("encrypted-notes"),
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockRepo, s.mockCrypto, s.logger)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
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
	cred, _ := NewCredential(testID, testName, testLogin, testPassword, testNotes, 1)
	return cred
}

func assertCredentialFields(t *testing.T, cred *Credential, userID string) {
	t.Helper()
	require.NotNil(t, cred)
	assert.Equal(t, testName, cred.GetName())
	assert.Equal(t, testLogin, cred.GetLogin())
	assert.Equal(t, testPassword, cred.GetPassword())
	assert.Equal(t, testNotes, cred.GetMetadata())
	assert.Equal(t, userID, cred.GetUserID())
}

func assertCredentialFieldsWithID(t *testing.T, cred *Credential, userID, credID string) {
	t.Helper()
	assertCredentialFields(t, cred, userID)
	assert.Equal(t, credID, cred.GetID())
}
