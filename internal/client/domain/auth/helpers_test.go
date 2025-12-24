package auth

import (
	"context"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
	grpcClient "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testLogin    = "testuser"
	testPassword = "testpass"
	testToken    = "token-xyz-456"
	testUserID   = "user-id-123"
)

type testSetup struct {
	ctrl        *gomock.Controller
	mockClient  *grpcClient.MockAuthClient
	mockStore   *MockSessionStore
	service     Service
	ctx         context.Context
	testSession shared.Session
	testCreds   interfaces.Credentials
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:        ctrl,
		mockClient:  grpcClient.NewMockAuthClient(ctrl),
		mockStore:   NewMockSessionStore(ctrl),
		ctx:         context.Background(),
		testSession: shared.NewSession(testToken, testToken, testUserID, testLogin),
		testCreds:   interfaces.NewCredentials(testLogin, testPassword),
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockClient, s.mockStore)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func (s *testSetup) expectLoginSuccess() {
	s.mockClient.EXPECT().
		Login(gomock.Any(), testLogin, testPassword).
		Return(s.testSession, nil).
		Times(1)
	s.mockStore.EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectLoginClientError(err error) {
	s.mockClient.EXPECT().
		Login(gomock.Any(), testLogin, testPassword).
		Return(shared.Session{}, err).
		Times(1)
}

func (s *testSetup) expectLoginSaveError(saveErr error) {
	s.mockClient.EXPECT().
		Login(gomock.Any(), testLogin, testPassword).
		Return(s.testSession, nil).
		Times(1)
	s.mockStore.EXPECT().
		Save(gomock.Any()).
		Return(saveErr).
		Times(1)
}

func (s *testSetup) expectRegisterSuccess() {
	s.mockClient.EXPECT().
		Register(gomock.Any(), testLogin, testPassword).
		Return(s.testSession, nil).
		Times(1)
	s.mockStore.EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectRegisterClientError(err error) {
	s.mockClient.EXPECT().
		Register(gomock.Any(), testLogin, testPassword).
		Return(shared.Session{}, err).
		Times(1)
}

func (s *testSetup) expectRegisterSaveError(saveErr error) {
	s.mockClient.EXPECT().
		Register(gomock.Any(), testLogin, testPassword).
		Return(s.testSession, nil).
		Times(1)
	s.mockStore.EXPECT().
		Save(gomock.Any()).
		Return(saveErr).
		Times(1)
}

func (s *testSetup) expectLogoutSuccess() {
	s.mockStore.EXPECT().
		Delete().
		Return(nil).
		Times(1)
}

func (s *testSetup) expectLogoutError(err error) {
	s.mockStore.EXPECT().
		Delete().
		Return(err).
		Times(1)
}

func (s *testSetup) expectGetSessionSuccess() {
	s.mockStore.EXPECT().
		Load().
		Return(s.testSession, nil).
		Times(1)
}

func (s *testSetup) expectGetSessionError(err error) {
	s.mockStore.EXPECT().
		Load().
		Return(shared.Session{}, err).
		Times(1)
}

func assertError(t *testing.T, err error, wantErr bool, expectedErr error) {
	t.Helper()
	if wantErr {
		require.Error(t, err)
		if expectedErr != nil {
			assert.ErrorIs(t, err, expectedErr)
		}
	} else {
		require.NoError(t, err)
	}
}

func assertSessionEqual(t *testing.T, got, want shared.Session) {
	t.Helper()
	assert.Equal(t, want, got)
}
