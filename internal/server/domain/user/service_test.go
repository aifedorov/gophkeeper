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
	"golang.org/x/crypto/bcrypt"
)

const (
	testLogin        = "testuser"
	testPass         = "testpassword"
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
				// Password hash is generated dynamically by bcrypt
				expectedUser := &repository2.User{
					ID:           expectedID,
					Login:        testLogin,
					PasswordHash: "",
				}
				m.EXPECT().
					CreateUser(testLogin, gomock.Any()).
					Times(1).
					DoAndReturn(func(login, passHash string) (*repository2.User, error) {
						expectedUser.PasswordHash = passHash
						return expectedUser, nil
					})
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
					CreateUser(existingLogin, gomock.Any()).
					Times(1).
					Return(nil, repository2.ErrLoginExists)
			},
			wantErrIs: ErrLoginExists,
		},
		{
			name:     "inMemory error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					CreateUser(testLogin, gomock.Any()).
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
					CreateUser("", gomock.Any()).
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
					CreateUser(testLogin, gomock.Any()).
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
				// Generate a bcrypt hash for the test password
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)
				expectedUser := &repository2.User{
					ID:           expectedID,
					Login:        testLogin,
					PasswordHash: string(hashedPassword),
				}
				m.EXPECT().
					GetUser(testLogin).
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
			name:     "auth.proto not found",
			login:    nonexistentLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					GetUser(nonexistentLogin).
					Times(1).
					Return(nil, repository2.ErrUserNotFound)
			},
			wantErrIs: ErrUserNotFound,
		},
		{
			name:     "inMemory error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, _ uuid.UUID) {
				m.EXPECT().
					GetUser(testLogin).
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
					GetUser("").
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
					GetUser(testLogin).
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

func TestLogin_InvalidPassword(t *testing.T) {
	t.Parallel()

	t.Run("returns error for invalid password", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		expectedID := uuid.New()

		// Generate hash for correct password
		correctPassword := "correctpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

		expectedUser := &repository2.User{
			ID:           expectedID,
			Login:        testLogin,
			PasswordHash: string(hashedPassword),
		}

		mockRepo.EXPECT().
			GetUser(testLogin).
			Times(1).
			Return(expectedUser, nil)

		logger := zap.NewNop()
		service := NewService(mockRepo, logger)

		// Try to login with wrong password
		user, err := service.Login(testLogin, "wrongpassword")

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to compare hash and password")
	})
}
