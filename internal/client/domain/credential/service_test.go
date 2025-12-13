package credential

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockCredentialClient(ctrl)
	service := NewService(mockClient)

	require.NotNil(t, service)
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cred      Credential
		setupMock func(*MockCredentialClient)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful creation",
			cred: Credential{
				ID:       "test-id-123",
				Name:     "test-credential",
				Login:    "testuser",
				Password: "testpass",
				Notes:    "test notes",
			},
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
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
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty name",
			cred: Credential{
				ID:       "test-id-789",
				Name:     "",
				Login:    "testuser",
				Password: "testpass",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "validation error - empty login",
			cred: Credential{
				ID:       "test-id-101",
				Name:     "test-name",
				Login:    "",
				Password: "testpass",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "validation error - empty password",
			cred: Credential{
				ID:       "test-id-102",
				Name:     "test-name",
				Login:    "testuser",
				Password: "",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
				// No expectation - validation fails before client call
			},
			wantErr: true,
			errMsg:  "invalid credential",
		},
		{
			name: "client error",
			cred: Credential{
				ID:       "test-id-103",
				Name:     "test-credential",
				Login:    "testuser",
				Password: "testpass",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("network error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to create credential",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockCredentialClient(ctrl)
			tt.setupMock(mockClient)

			service := NewService(mockClient)
			err := service.Create(context.Background(), tt.cred)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		setupMock func(*MockCredentialClient)
		wantCred  Credential
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful get",
			id:   "test-id-123",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Get(gomock.Any(), "test-id-123").
					Return(Credential{
						ID:       "test-id-123",
						Name:     "test-credential",
						Login:    "testuser",
						Password: "testpass",
						Notes:    "test notes",
					}, nil).
					Times(1)
			},
			wantCred: Credential{
				ID:       "test-id-123",
				Name:     "test-credential",
				Login:    "testuser",
				Password: "testpass",
				Notes:    "test notes",
			},
			wantErr: false,
		},
		{
			name: "successful get without notes",
			id:   "test-id-456",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Get(gomock.Any(), "test-id-456").
					Return(Credential{
						ID:       "test-id-456",
						Name:     "another-credential",
						Login:    "anotheruser",
						Password: "anotherpass",
						Notes:    "",
					}, nil).
					Times(1)
			},
			wantCred: Credential{
				ID:       "test-id-456",
				Name:     "another-credential",
				Login:    "anotheruser",
				Password: "anotherpass",
				Notes:    "",
			},
			wantErr: false,
		},
		{
			name: "client error - not found",
			id:   "non-existent-id",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Get(gomock.Any(), "non-existent-id").
					Return(Credential{}, errors.New("not found")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get credential",
		},
		{
			name: "client error - network error",
			id:   "test-id-789",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Get(gomock.Any(), "test-id-789").
					Return(Credential{}, errors.New("network error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get credential",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockCredentialClient(ctrl)
			tt.setupMock(mockClient)

			service := NewService(mockClient)
			cred, err := service.Get(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCred, cred)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*MockCredentialClient)
		wantCreds []Credential
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful list with multiple credentials",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					List(gomock.Any()).
					Return([]Credential{
						{
							ID:       "id-1",
							Name:     "cred-1",
							Login:    "user1",
							Password: "pass1",
							Notes:    "notes1",
						},
						{
							ID:       "id-2",
							Name:     "cred-2",
							Login:    "user2",
							Password: "pass2",
							Notes:    "notes2",
						},
					}, nil).
					Times(1)
			},
			wantCreds: []Credential{
				{
					ID:       "id-1",
					Name:     "cred-1",
					Login:    "user1",
					Password: "pass1",
					Notes:    "notes1",
				},
				{
					ID:       "id-2",
					Name:     "cred-2",
					Login:    "user2",
					Password: "pass2",
					Notes:    "notes2",
				},
			},
			wantErr: false,
		},
		{
			name: "successful list with empty result",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					List(gomock.Any()).
					Return([]Credential{}, nil).
					Times(1)
			},
			wantCreds: []Credential{},
			wantErr:   false,
		},
		{
			name: "successful list with single credential",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					List(gomock.Any()).
					Return([]Credential{
						{
							ID:       "single-id",
							Name:     "single-cred",
							Login:    "singleuser",
							Password: "singlepass",
							Notes:    "",
						},
					}, nil).
					Times(1)
			},
			wantCreds: []Credential{
				{
					ID:       "single-id",
					Name:     "single-cred",
					Login:    "singleuser",
					Password: "singlepass",
					Notes:    "",
				},
			},
			wantErr: false,
		},
		{
			name: "client error",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("network error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get list of credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockCredentialClient(ctrl)
			tt.setupMock(mockClient)

			service := NewService(mockClient)
			creds, err := service.List(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCreds, creds)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		cred      Credential
		setupMock func(*MockCredentialClient)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful update",
			id:   "test-id-123",
			cred: Credential{
				ID:       "test-id-123",
				Name:     "updated-credential",
				Login:    "updateduser",
				Password: "updatedpass",
				Notes:    "updated notes",
			},
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Update(gomock.Any(), "test-id-123", gomock.Any()).
					Return(nil).
					Times(1)
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
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Update(gomock.Any(), "test-id-456", gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty name",
			id:   "test-id-789",
			cred: Credential{
				ID:       "test-id-789",
				Name:     "",
				Login:    "testuser",
				Password: "testpass",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
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
				Name:     "test-name",
				Login:    "",
				Password: "testpass",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
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
				Name:     "test-name",
				Login:    "testuser",
				Password: "",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
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
				Name:     "test-credential",
				Login:    "testuser",
				Password: "testpass",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Update(gomock.Any(), "non-existent-id", gomock.Any()).
					Return(errors.New("not found")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get credential",
		},
		{
			name: "client error - network error",
			id:   "test-id-103",
			cred: Credential{
				ID:       "test-id-103",
				Name:     "test-credential",
				Login:    "testuser",
				Password: "testpass",
				Notes:    "notes",
			},
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Update(gomock.Any(), "test-id-103", gomock.Any()).
					Return(errors.New("network error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get credential",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockCredentialClient(ctrl)
			tt.setupMock(mockClient)

			service := NewService(mockClient)
			err := service.Update(context.Background(), tt.id, tt.cred)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		setupMock func(*MockCredentialClient)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful deletion",
			id:   "test-id-123",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Delete(gomock.Any(), "test-id-123").
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "successful deletion with different ID",
			id:   "another-id-456",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Delete(gomock.Any(), "another-id-456").
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "client error - not found",
			id:   "non-existent-id",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Delete(gomock.Any(), "non-existent-id").
					Return(errors.New("not found")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to delete credential",
		},
		{
			name: "client error - network error",
			id:   "test-id-789",
			setupMock: func(m *MockCredentialClient) {
				m.EXPECT().
					Delete(gomock.Any(), "test-id-789").
					Return(errors.New("network error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to delete credential",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockCredentialClient(ctrl)
			tt.setupMock(mockClient)

			service := NewService(mockClient)
			err := service.Delete(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
