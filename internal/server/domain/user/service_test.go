package user

import (
	"errors"
	"testing"

	repository2 "github.com/aifedorov/gophkeeper/internal/server/domain/user/repository/db"
	"github.com/aifedorov/gophkeeper/internal/server/domain/user/repository/db/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	testLogin        = "testuser"
	testPass         = "hashedpassword"
	existingLogin    = "existinguser"
	nonexistentLogin = "nonexistentuser"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		login        string
		passHash     string
		setupMock    func(*mocks.MockRepository, uuid.UUID)
		wantErr      error
		wantErrIs    error
		validateUser func(*testing.T, *User, uuid.UUID)
	}{
		{
			name:     "successful registration",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, expectedID uuid.UUID) {
				expectedUser := &repository2.User{
					ID:           expectedID,
					Login:        testLogin,
					PasswordHash: testPass,
				}
				m.EXPECT().
					CreateUser(testLogin, testPass).
					Times(1).
					Return(expectedUser, nil)
			},
			validateUser: func(t *testing.T, user *User, expectedID uuid.UUID) {
				require.NotNil(t, user)
				assert.Equal(t, testLogin, user.GetLogin())
				assert.Equal(t, expectedID.String(), user.GetUserID())
			},
		},
		{
			name:     "login already exists",
			login:    existingLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					CreateUser(existingLogin, testPass).
					Times(1).
					Return(nil, repository2.ErrLoginExists)
			},
			wantErrIs: ErrLoginExists,
		},
		{
			name:     "repository error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					CreateUser(testLogin, testPass).
					Times(1).
					Return(nil, errors.New("database connection failed"))
			},
			wantErr: errors.New("database connection failed"),
		},
		{
			name:     "empty login",
			login:    "",
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					CreateUser("", testPass).
					Times(1).
					Return(nil, errors.New("invalid login"))
			},
			wantErr: errors.New("invalid login"),
		},
		{
			name:     "empty password hash",
			login:    testLogin,
			passHash: "",
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					CreateUser(testLogin, "").
					Times(1).
					Return(nil, errors.New("invalid password"))
			},
			wantErr: errors.New("invalid password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			expectedID := uuid.New()
			tt.setupMock(mockRepo, expectedID)

			logger := zap.NewNop()
			service := NewService(mockRepo, logger)

			user, err := service.Register(tt.login, tt.passHash)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, user)
			} else if tt.wantErr != nil {
				require.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				tt.validateUser(t, user, expectedID)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		login        string
		passHash     string
		setupMock    func(*mocks.MockRepository, uuid.UUID)
		wantErr      error
		wantErrIs    error
		validateUser func(*testing.T, *User, uuid.UUID)
	}{
		{
			name:     "successful login",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, expectedID uuid.UUID) {
				expectedUser := &repository2.User{
					ID:           expectedID,
					Login:        testLogin,
					PasswordHash: testPass,
				}
				m.EXPECT().
					GetUser(testLogin, testPass).
					Times(1).
					Return(expectedUser, nil)
			},
			validateUser: func(t *testing.T, user *User, expectedID uuid.UUID) {
				require.NotNil(t, user)
				assert.Equal(t, testLogin, user.GetLogin())
				assert.Equal(t, expectedID.String(), user.GetUserID())
			},
		},
		{
			name:     "user not found",
			login:    nonexistentLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					GetUser(nonexistentLogin, testPass).
					Times(1).
					Return(nil, repository2.ErrUserNotFound)
			},
			wantErrIs: ErrUserNotFound,
		},
		{
			name:     "repository error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					GetUser(testLogin, testPass).
					Times(1).
					Return(nil, errors.New("database connection failed"))
			},
			wantErr: errors.New("database connection failed"),
		},
		{
			name:     "empty login",
			login:    "",
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					GetUser("", testPass).
					Times(1).
					Return(nil, repository2.ErrUserNotFound)
			},
			wantErrIs: ErrUserNotFound,
		},
		{
			name:     "empty password hash",
			login:    testLogin,
			passHash: "",
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					GetUser(testLogin, "").
					Times(1).
					Return(nil, repository2.ErrUserNotFound)
			},
			wantErrIs: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			expectedID := uuid.New()
			tt.setupMock(mockRepo, expectedID)

			logger := zap.NewNop()
			service := NewService(mockRepo, logger)

			user, err := service.Login(tt.login, tt.passHash)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, user)
			} else if tt.wantErr != nil {
				require.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				tt.validateUser(t, user, expectedID)
			}
		})
	}
}
