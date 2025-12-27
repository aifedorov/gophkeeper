package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
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
		setupMock    func(*mocks.MockRepository, *mocks.MockCryptoService, uuid.UUID)
		wantErr      error
		wantErrIs    error
		validateUser func(*testing.T, *User, uuid.UUID)
	}{
		{
			name:     "successful registration",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, expectedID uuid.UUID) {
				testSalt := []byte("test-salt-32-bytes-long-string!!")
				encryptionKey := argon2.IDKey([]byte(testPass), testSalt, 1, 64*1024, 4, 32)
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)

				c.EXPECT().GenerateSalt().Return(testSalt, nil).Times(1)
				c.EXPECT().DeriveEncryptionKey(testPass, string(testSalt)).Return(encryptionKey).Times(1)
				c.EXPECT().HashPassword(testPass).Return(string(hashedPassword), nil).Times(1)

				expectedUser := interfaces.RepositoryUser{
					ID:           expectedID.String(),
					Login:        testLogin,
					PasswordHash: string(hashedPassword),
					Salt:         string(testSalt),
				}
				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, user interfaces.RepositoryUser, passHash string) (interfaces.RepositoryUser, error) {
						expectedUser.PasswordHash = passHash
						expectedUser.Salt = user.Salt
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
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				testSalt := []byte("test-salt-32-bytes-long-string!!")
				encryptionKey := argon2.IDKey([]byte(testPass), testSalt, 1, 64*1024, 4, 32)
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)

				c.EXPECT().GenerateSalt().Return(testSalt, nil).Times(1)
				c.EXPECT().DeriveEncryptionKey(testPass, string(testSalt)).Return(encryptionKey).Times(1)
				c.EXPECT().HashPassword(testPass).Return(string(hashedPassword), nil).Times(1)

				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(interfaces.RepositoryUser{}, ErrLoginExists)
			},
			wantErrIs: ErrLoginExists,
		},
		{
			name:     "inMemory error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				testSalt := []byte("test-salt-32-bytes-long-string!!")
				encryptionKey := argon2.IDKey([]byte(testPass), testSalt, 1, 64*1024, 4, 32)
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)

				c.EXPECT().GenerateSalt().Return(testSalt, nil).Times(1)
				c.EXPECT().DeriveEncryptionKey(testPass, string(testSalt)).Return(encryptionKey).Times(1)
				c.EXPECT().HashPassword(testPass).Return(string(hashedPassword), nil).Times(1)

				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(interfaces.RepositoryUser{}, errors.New("database connection failed"))
			},
			wantErr: errors.New("database connection failed"),
		},
		{
			name:     "empty login",
			login:    "",
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				// No mock expectation - validation fails before repository call
			},
			wantErr: errors.New("invalid login"),
		},
		{
			name:     "empty password hash",
			login:    testLogin,
			passHash: "",
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				// No mock expectation - validation fails before repository call
			},
			wantErr: errors.New("invalid password"),
		},
		{
			name:     "salt generation error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				c.EXPECT().GenerateSalt().Return(nil, errors.New("entropy failure")).Times(1)
			},
			wantErr: errors.New("failed to generate salt"),
		},
		{
			name:     "password hash error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				testSalt := []byte("test-salt-32-bytes-long-string!!")
				encryptionKey := argon2.IDKey([]byte(testPass), testSalt, 1, 64*1024, 4, 32)

				c.EXPECT().GenerateSalt().Return(testSalt, nil).Times(1)
				c.EXPECT().DeriveEncryptionKey(testPass, string(testSalt)).Return(encryptionKey).Times(1)
				c.EXPECT().HashPassword(testPass).Return("", errors.New("hash failure")).Times(1)
			},
			wantErr: errors.New("failed to hash password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			mockCrypto := mocks.NewMockCryptoService(ctrl)
			expectedID := uuid.New()
			tt.setupMock(mockRepo, mockCrypto, expectedID)

			logger := zap.NewNop()
			service := NewService(mockRepo, logger, mockCrypto)

			ctx := context.Background()
			user, _, err := service.Register(ctx, tt.login, tt.passHash)

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
		setupMock    func(*mocks.MockRepository, *mocks.MockCryptoService, uuid.UUID)
		wantErr      error
		wantErrIs    error
		validateUser func(*testing.T, *User, uuid.UUID)
	}{
		{
			name:     "successful login",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, expectedID uuid.UUID) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)
				testSalt := "testsalt"
				encryptionKey := argon2.IDKey([]byte(testPass), []byte(testSalt), 1, 64*1024, 4, 32)

				expectedUser := interfaces.RepositoryUser{
					ID:           expectedID.String(),
					Login:        testLogin,
					PasswordHash: string(hashedPassword),
					Salt:         testSalt,
				}
				m.EXPECT().
					GetUser(gomock.Any(), testLogin).
					Times(1).
					Return(expectedUser, nil)
				c.EXPECT().CompareHashAndPassword(string(hashedPassword), testPass).Return(nil).Times(1)
				c.EXPECT().DeriveEncryptionKey(testPass, testSalt).Return(encryptionKey).Times(1)
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
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				m.EXPECT().
					GetUser(gomock.Any(), nonexistentLogin).
					Times(1).
					Return(interfaces.RepositoryUser{}, ErrUserNotFound)
			},
			wantErrIs: ErrUserNotFound,
		},
		{
			name:     "inMemory error",
			login:    testLogin,
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				m.EXPECT().
					GetUser(gomock.Any(), testLogin).
					Times(1).
					Return(interfaces.RepositoryUser{}, errors.New("database connection failed"))
			},
			wantErr: errors.New("database connection failed"),
		},
		{
			name:     "empty login",
			login:    "",
			passHash: testPass,
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				// No mock expectation - validation fails before repository call
			},
			wantErr: errors.New("invalid login"),
		},
		{
			name:     "empty password hash",
			login:    testLogin,
			passHash: "",
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, _ uuid.UUID) {
				// No mock expectation - validation fails before repository call
			},
			wantErr: errors.New("invalid password"),
		},
		{
			name:     "invalid credentials - password mismatch",
			login:    testLogin,
			passHash: "wrongpassword",
			setupMock: func(m *mocks.MockRepository, c *mocks.MockCryptoService, expectedID uuid.UUID) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)
				testSalt := "testsalt"
				encryptionKey := argon2.IDKey([]byte("wrongpassword"), []byte(testSalt), 1, 64*1024, 4, 32)

				expectedUser := interfaces.RepositoryUser{
					ID:           expectedID.String(),
					Login:        testLogin,
					PasswordHash: string(hashedPassword),
					Salt:         testSalt,
				}
				m.EXPECT().
					GetUser(gomock.Any(), testLogin).
					Times(1).
					Return(expectedUser, nil)
				c.EXPECT().DeriveEncryptionKey("wrongpassword", testSalt).Return(encryptionKey).Times(1)
				c.EXPECT().CompareHashAndPassword(string(hashedPassword), "wrongpassword").Return(errors.New("password mismatch")).Times(1)
			},
			wantErr: errors.New("failed to compare hash and password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			mockCrypto := mocks.NewMockCryptoService(ctrl)
			expectedID := uuid.New()
			tt.setupMock(mockRepo, mockCrypto, expectedID)

			logger := zap.NewNop()
			service := NewService(mockRepo, logger, mockCrypto)

			ctx := context.Background()
			user, _, err := service.Login(ctx, tt.login, tt.passHash)

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

		correctPassword := "correctpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

		expectedUser := interfaces.RepositoryUser{
			ID:           expectedID.String(),
			Login:        testLogin,
			PasswordHash: string(hashedPassword),
			Salt:         "testsalt",
		}

		mockRepo.EXPECT().
			GetUser(gomock.Any(), testLogin).
			Times(1).
			Return(expectedUser, nil)

		mockCrypto := mocks.NewMockCryptoService(ctrl)
		mockCrypto.EXPECT().DeriveEncryptionKey("wrongpassword", "testsalt").
			Return([]byte("derived-key")).Times(1)
		mockCrypto.EXPECT().CompareHashAndPassword(string(hashedPassword), "wrongpassword").
			Return(bcrypt.ErrMismatchedHashAndPassword).Times(1)

		logger := zap.NewNop()
		service := NewService(mockRepo, logger, mockCrypto)

		ctx := context.Background()
		user, _, err := service.Login(ctx, testLogin, "wrongpassword")

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to compare hash and password")
	})
}

func TestService_GetUserDataFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupContext func(*testing.T, Service) context.Context
		wantUserID   string
		wantEncKey   string
		wantErr      bool
		errMsg       string
	}{
		{
			name: "successful get user data",
			setupContext: func(t *testing.T, service Service) context.Context {
				ctx := context.Background()
				ctx = service.SetUserID(ctx, "test-user-123")
				ctx = service.SetEncryptionKeyEncoded(ctx, "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEhAAA=")
				return ctx
			},
			wantUserID: "test-user-123",
			wantEncKey: "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEhAAA=",
			wantErr:    false,
		},
		{
			name: "missing user ID",
			setupContext: func(t *testing.T, service Service) context.Context {
				ctx := context.Background()
				ctx = service.SetEncryptionKeyEncoded(ctx, "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEhAAA=")
				return ctx
			},
			wantErr: true,
			errMsg:  "failed to get userID",
		},
		{
			name: "missing encryption key",
			setupContext: func(t *testing.T, service Service) context.Context {
				ctx := context.Background()
				ctx = service.SetUserID(ctx, "test-user-123")
				return ctx
			},
			wantErr: true,
			errMsg:  "failed to get encryption key",
		},
		{
			name: "both missing",
			setupContext: func(t *testing.T, service Service) context.Context {
				return context.Background()
			},
			wantErr: true,
			errMsg:  "failed to get userID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			mockCrypto := mocks.NewMockCryptoService(ctrl)
			logger := zap.NewNop()
			service := NewService(mockRepo, logger, mockCrypto)

			ctx := tt.setupContext(t, service)
			userID, encKey, err := service.GetUserDataFromContext(ctx)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Empty(t, userID)
				assert.Empty(t, encKey)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUserID, userID)
				assert.Equal(t, tt.wantEncKey, encKey)
			}
		})
	}
}
