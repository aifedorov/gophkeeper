package credential

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testID       = "test-id-123"
	testName     = "test-credential"
	testLogin    = "testuser"
	testPassword = "testpass"
	testNotes    = "test notes"
	testVersion  = int64(1)
)

type testSetup struct {
	ctrl       *gomock.Controller
	mockClient *MockCredentialClient
	mockCache  *MockCacheStorage
	service    Service
	ctx        context.Context
	testCred   Credential
	wantCreds  []Credential
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:       ctrl,
		mockClient: NewMockCredentialClient(ctrl),
		mockCache:  NewMockCacheStorage(ctrl),
		ctx:        context.Background(),
		testCred: Credential{
			ID:       testID,
			Name:     testName,
			Login:    testLogin,
			Password: testPassword,
			Notes:    testNotes,
		},
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockClient, s.mockCache)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func (s *testSetup) expectCreateSuccess(id string, version int64) {
	s.mockClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(id, version, nil).
		Times(1)
	s.mockCache.EXPECT().
		SetCredentialVersion(id, version).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectCreateClientError(err error) {
	s.mockClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return("", int64(0), err).
		Times(1)
}

func (s *testSetup) expectListSuccess(creds []Credential) {
	s.mockClient.EXPECT().
		List(gomock.Any()).
		Return(creds, nil).
		Times(1)
	for _, cred := range creds {
		if cred.Version > 0 {
			s.mockCache.EXPECT().
				SetCredentialVersion(cred.ID, cred.Version).
				Return(nil).
				Times(1)
		}
	}
}

func (s *testSetup) expectListError(err error) {
	s.mockClient.EXPECT().
		List(gomock.Any()).
		Return(nil, err).
		Times(1)
}

func (s *testSetup) expectUpdateSuccess(id string, currentVersion, newVersion int64) {
	s.mockCache.EXPECT().
		GetCredentialVersion(id).
		Return(currentVersion, nil).
		Times(1)
	s.mockClient.EXPECT().
		Update(gomock.Any(), id, gomock.Any()).
		Return(newVersion, nil).
		Times(1)
	s.mockCache.EXPECT().
		SetCredentialVersion(id, newVersion).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectUpdateClientError(id string, currentVersion int64, err error) {
	s.mockCache.EXPECT().
		GetCredentialVersion(id).
		Return(currentVersion, nil).
		Times(1)
	s.mockClient.EXPECT().
		Update(gomock.Any(), id, gomock.Any()).
		Return(int64(0), err).
		Times(1)
}

func (s *testSetup) expectDeleteSuccess(id string) {
	s.mockClient.EXPECT().
		Delete(gomock.Any(), id).
		Return(nil).
		Times(1)
	s.mockCache.EXPECT().
		DeleteCredentialVersion(id).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectDeleteClientError(id string, err error) {
	s.mockClient.EXPECT().
		Delete(gomock.Any(), id).
		Return(err).
		Times(1)
}

func assertError(t *testing.T, err error, wantErr bool, errMsg string) {
	t.Helper()
	if wantErr {
		require.Error(t, err)
		if errMsg != "" {
			assert.Contains(t, err.Error(), errMsg)
		}
	} else {
		require.NoError(t, err)
	}
}

func assertCredsEqual(t *testing.T, got, want []Credential) {
	t.Helper()
	assert.Equal(t, want, got)
}
