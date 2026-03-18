package credential

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	s := newTestSetup(t)
	defer s.cleanup()
	service := NewService(s.mockClient, s.mockCache)

	require.NotNil(t, service)
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cred    Credential
		setup   func(*testSetup, Credential)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful creation",
			cred: Credential{
				ID:       testID,
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, cred Credential) {
				s.expectCreateSuccess(cred.ID, testVersion)
			},
			wantErr: false,
		},
		{
			name: "successful creation without notes",
			cred: Credential{
				ID:       "test-id-456",
				Name:     "another-credential",
				Login:    "anotheruser",
				Password: "anotherpass",
				Notes:    "",
			},
			setup: func(s *testSetup, cred Credential) {
				s.expectCreateSuccess(cred.ID, testVersion)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty name",
			cred: Credential{
				ID:       "test-id-789",
				Name:     "",
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, cred Credential) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "validation error - empty login",
			cred: Credential{
				ID:       "test-id-101",
				Name:     testName,
				Login:    "",
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, cred Credential) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "validation error - empty password",
			cred: Credential{
				ID:       "test-id-102",
				Name:     testName,
				Login:    testLogin,
				Password: "",
				Notes:    testNotes,
			},
			setup: func(s *testSetup, cred Credential) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "client error",
			cred: Credential{
				ID:       "test-id-103",
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, cred Credential) {
				s.expectCreateClientError(errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to create credential",
		},
		{
			name: "cache error on set version",
			cred: Credential{
				ID:       "cache-error-id",
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, cred Credential) {
				s.mockClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return("cache-error-id", int64(1), nil).
					Times(1)
				s.mockCache.EXPECT().
					SetCredentialVersion("cache-error-id", int64(1)).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to save credential to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.cred)

			err := s.service.Create(s.ctx, tt.cred)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*testSetup)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful list with multiple credentials",
			setup: func(s *testSetup) {
				creds := []Credential{
					{ID: "id-1", Name: "cred-1", Login: "user1", Password: "pass1", Notes: "notes1", Version: 1},
					{ID: "id-2", Name: "cred-2", Login: "user2", Password: "pass2", Notes: "notes2", Version: 1},
				}
				s.expectListSuccess(creds)
				s.wantCreds = creds
			},
			wantErr: false,
		},
		{
			name: "successful list with empty result",
			setup: func(s *testSetup) {
				s.expectListSuccess([]Credential{})
				s.wantCreds = []Credential{}
			},
			wantErr: false,
		},
		{
			name: "successful list with single credential",
			setup: func(s *testSetup) {
				creds := []Credential{
					{ID: "single-id", Name: "single-cred", Login: "singleuser", Password: "singlepass", Notes: "", Version: 2},
				}
				s.expectListSuccess(creds)
				s.wantCreds = creds
			},
			wantErr: false,
		},
		{
			name: "client error",
			setup: func(s *testSetup) {
				s.expectListError(errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to get list of credentials",
		},
		{
			name: "server returns invalid version 0",
			setup: func(s *testSetup) {
				creds := []Credential{
					{ID: "id-1", Name: "cred-1", Login: "user1", Password: "pass1", Notes: "", Version: 0},
				}
				s.mockClient.EXPECT().
					List(gomock.Any()).
					Return(creds, nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "server returned invalid version 0",
		},
		{
			name: "cache error on set version",
			setup: func(s *testSetup) {
				creds := []Credential{
					{ID: "id-1", Name: "cred-1", Login: "user1", Password: "pass1", Notes: "", Version: 1},
				}
				s.mockClient.EXPECT().
					List(gomock.Any()).
					Return(creds, nil).
					Times(1)
				s.mockCache.EXPECT().
					SetCredentialVersion("id-1", int64(1)).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to save credential to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			creds, err := s.service.List(s.ctx)

			assertError(t, err, tt.wantErr, tt.errMsg)
			if !tt.wantErr && s.wantCreds != nil {
				assertCredsEqual(t, creds, s.wantCreds)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		cred    Credential
		setup   func(*testSetup, string, Credential)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful update",
			id:   testID,
			cred: Credential{
				ID:       testID,
				Name:     "updated-credential",
				Login:    "updateduser",
				Password: "updatedpass",
				Notes:    "updated notes",
			},
			setup: func(s *testSetup, id string, cred Credential) {
				s.expectUpdateSuccess(id, testVersion, int64(2))
			},
			wantErr: false,
		},
		{
			name: "successful update without notes",
			id:   "test-id-456",
			cred: Credential{
				ID:       "test-id-456",
				Name:     "updated-credential",
				Login:    "updateduser",
				Password: "updatedpass",
				Notes:    "",
			},
			setup: func(s *testSetup, id string, cred Credential) {
				s.expectUpdateSuccess(id, int64(3), int64(4))
			},
			wantErr: false,
		},
		{
			name: "validation error - empty name",
			id:   "test-id-789",
			cred: Credential{
				ID:       "test-id-789",
				Name:     "",
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, id string, cred Credential) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "validation error - empty login",
			id:   "test-id-101",
			cred: Credential{
				ID:       "test-id-101",
				Name:     testName,
				Login:    "",
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, id string, cred Credential) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "validation error - empty password",
			id:   "test-id-102",
			cred: Credential{
				ID:       "test-id-102",
				Name:     testName,
				Login:    testLogin,
				Password: "",
				Notes:    testNotes,
			},
			setup: func(s *testSetup, id string, cred Credential) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "client error - not found",
			id:   "non-existent-id",
			cred: Credential{
				ID:       "non-existent-id",
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, id string, cred Credential) {
				s.expectUpdateClientError(id, testVersion, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to update",
		},
		{
			name: "client error - network error",
			id:   "test-id-103",
			cred: Credential{
				ID:       "test-id-103",
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, id string, cred Credential) {
				s.expectUpdateClientError(id, int64(2), errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to update",
		},
		{
			name: "cache error - get version fails",
			id:   "test-id-104",
			cred: Credential{
				ID:       "test-id-104",
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, id string, cred Credential) {
				s.mockCache.EXPECT().
					GetCredentialVersion(id).
					Return(int64(0), errors.New("cache miss")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get version from cache",
		},
		{
			name: "cache error - set version fails after update",
			id:   "test-id-105",
			cred: Credential{
				ID:       "test-id-105",
				Name:     testName,
				Login:    testLogin,
				Password: testPassword,
				Notes:    testNotes,
			},
			setup: func(s *testSetup, id string, cred Credential) {
				s.mockCache.EXPECT().
					GetCredentialVersion(id).
					Return(int64(1), nil).
					Times(1)
				s.mockClient.EXPECT().
					Update(gomock.Any(), id, gomock.Any()).
					Return(int64(2), nil).
					Times(1)
				s.mockCache.EXPECT().
					SetCredentialVersion(gomock.Any(), int64(2)).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to save credential to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.id, tt.cred)

			err := s.service.Update(s.ctx, tt.id, tt.cred)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		setup   func(*testSetup, string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful deletion",
			id:   testID,
			setup: func(s *testSetup, id string) {
				s.expectDeleteSuccess(id)
			},
			wantErr: false,
		},
		{
			name: "successful deletion with different ID",
			id:   "another-id-456",
			setup: func(s *testSetup, id string) {
				s.expectDeleteSuccess(id)
			},
			wantErr: false,
		},
		{
			name: "client error - not found",
			id:   "non-existent-id",
			setup: func(s *testSetup, id string) {
				s.expectDeleteClientError(id, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to delete credential",
		},
		{
			name: "client error - network error",
			id:   "test-id-789",
			setup: func(s *testSetup, id string) {
				s.expectDeleteClientError(id, errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to delete credential",
		},
		{
			name: "cache error - delete version fails",
			id:   "test-id-999",
			setup: func(s *testSetup, id string) {
				s.mockClient.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil).
					Times(1)
				s.mockCache.EXPECT().
					DeleteCredentialVersion(id).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to delete credential from cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.id)

			err := s.service.Delete(s.ctx, tt.id)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}
