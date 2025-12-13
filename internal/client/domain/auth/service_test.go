package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	grpcClient "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testLogin    = "testuser"
	testPassword = "testpass"
	testToken    = "token-xyz-456"
)

func TestService_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		creds     interfaces.Credentials
		setupMock func(*grpcClient.MockAuthClient, *MockRepository)
		wantErr   bool
		errCheck  func(*testing.T, error)
	}{
		{
			name:  "successful login",
			creds: interfaces.NewCredentials(testLogin, testPassword),
			setupMock: func(client *grpcClient.MockAuthClient, repo *MockRepository) {
				client.EXPECT().
					Login(gomock.Any(), testLogin, testPassword).
					Return(testToken, []byte(testToken), nil)

				repo.EXPECT().
					Save(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "login fails - client error",
			creds: interfaces.NewCredentials(testLogin, testPassword),
			setupMock: func(client *grpcClient.MockAuthClient, repo *MockRepository) {
				client.EXPECT().
					Login(gomock.Any(), testLogin, testPassword).
					Return("", nil, ErrInvalidCredentials)
			},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		{
			name:  "login succeeds but save fails",
			creds: interfaces.NewCredentials(testLogin, testPassword),
			setupMock: func(client *grpcClient.MockAuthClient, repo *MockRepository) {
				client.EXPECT().
					Login(gomock.Any(), testLogin, testPassword).
					Return(testToken, []byte(testToken), nil)

				repo.EXPECT().
					Save(gomock.Any()).
					Return(errors.New("storage error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := grpcClient.NewMockAuthClient(ctrl)
			mockRepo := NewMockRepository(ctrl)
			tt.setupMock(mockClient, mockRepo)

			service := NewService(mockClient, mockRepo)

			err := service.Login(context.Background(), tt.creds)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errCheck != nil {
					tt.errCheck(t, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		creds     interfaces.Credentials
		setupMock func(*grpcClient.MockAuthClient, *MockRepository)
		wantErr   bool
		errCheck  func(*testing.T, error)
	}{
		{
			name:  "successful registration",
			creds: interfaces.NewCredentials(testLogin, testPassword),
			setupMock: func(client *grpcClient.MockAuthClient, repo *MockRepository) {
				client.EXPECT().
					Register(gomock.Any(), testLogin, testPassword).
					Return(testToken, []byte(testToken), nil)

				repo.EXPECT().
					Save(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "registration fails - auth already exists",
			creds: interfaces.NewCredentials(testLogin, testPassword),
			setupMock: func(client *grpcClient.MockAuthClient, repo *MockRepository) {
				client.EXPECT().
					Register(gomock.Any(), testLogin, testPassword).
					Return("", nil, ErrUserAlreadyExists)
			},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrUserAlreadyExists)
			},
		},
		{
			name:  "registration succeeds but save fails",
			creds: interfaces.NewCredentials(testLogin, testPassword),
			setupMock: func(client *grpcClient.MockAuthClient, repo *MockRepository) {
				client.EXPECT().
					Register(gomock.Any(), testLogin, testPassword).
					Return(testToken, []byte(testToken), nil)

				repo.EXPECT().
					Save(gomock.Any()).
					Return(errors.New("storage error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := grpcClient.NewMockAuthClient(ctrl)
			mockRepo := NewMockRepository(ctrl)
			tt.setupMock(mockClient, mockRepo)

			service := NewService(mockClient, mockRepo)

			err := service.Register(context.Background(), tt.creds)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errCheck != nil {
					tt.errCheck(t, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Logout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*MockRepository)
		wantErr   bool
	}{
		{
			name: "successful logout",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete().
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "logout fails - delete error",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete().
					Return(errors.New("delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := grpcClient.NewMockAuthClient(ctrl)
			mockRepo := NewMockRepository(ctrl)
			tt.setupMock(mockRepo)

			service := NewService(mockClient, mockRepo)

			err := service.Logout(context.Background())

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetCurrentSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(*MockRepository)
		wantSession interfaces.Session
		wantErr     bool
		errCheck    func(*testing.T, error)
	}{
		{
			name: "get session successfully",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().
					Load().
					Return(interfaces.NewSession(testToken, testToken), nil)
			},
			wantSession: interfaces.NewSession(testToken, testToken),
			wantErr:     false,
		},
		{
			name: "session not found",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().
					Load().
					Return(interfaces.Session{}, ErrSessionNotFound)
			},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrSessionNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := grpcClient.NewMockAuthClient(ctrl)
			mockRepo := NewMockRepository(ctrl)
			tt.setupMock(mockRepo)

			service := NewService(mockClient, mockRepo)

			session, err := service.GetCurrentSession()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errCheck != nil {
					tt.errCheck(t, err)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantSession, session)
			}
		})
	}
}
