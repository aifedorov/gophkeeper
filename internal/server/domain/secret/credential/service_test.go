package credential

import (
	"context"
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "successful creation",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				repoCred := interfaces.RepositoryCredential{
					ID:                s.credentialID,
					UserID:            s.userID,
					Name:              testName,
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					CreateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCred, nil).
					Times(1)

				s.expectDecryptCredential()
			},
		},
		{
			name: "name already exists",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				s.mockRepo.EXPECT().
					CreateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, ErrNameExists).
					Times(1)
			},
			wantErrIs: ErrNameExists,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				s.mockRepo.EXPECT().
					CreateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "repository returns nil",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				s.mockRepo.EXPECT().
					CreateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "encryption fails on login",
			setupMock: func(s *testSetup) {
				s.mockCrypto.EXPECT().
					Encrypt(testLogin, s.encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "decryption fails after creation",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				repoCred := interfaces.RepositoryCredential{
					ID:                s.credentialID,
					UserID:            s.userID,
					Name:              testName,
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					CreateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCred, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedLogin, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			cred := newTestCredential()
			result, err := setup.service.Create(ctx, setup.userID, setup.encryptionKeyStr, *cred)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, result)
			} else if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assertCredentialFieldsWithID(t, result, setup.userID, setup.credentialID)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "successful retrieval",
			setupMock: func(s *testSetup) {
				repoCred := interfaces.RepositoryCredential{
					ID:                s.credentialID,
					UserID:            s.userID,
					Name:              testName,
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					GetCredential(gomock.Any(), s.userID, s.credentialID).
					Return(&repoCred, nil).
					Times(1)

				s.expectDecryptCredential()
			},
		},
		{
			name: "credential not found",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					GetCredential(gomock.Any(), s.userID, s.credentialID).
					Return(nil, ErrNotFound).
					Times(1)
			},
			wantErrIs: ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					GetCredential(gomock.Any(), s.userID, s.credentialID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "repository returns nil",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					GetCredential(gomock.Any(), s.userID, s.credentialID).
					Return(nil, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "decryption fails",
			setupMock: func(s *testSetup) {
				repoCred := interfaces.RepositoryCredential{
					ID:                s.credentialID,
					UserID:            s.userID,
					Name:              testName,
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					GetCredential(gomock.Any(), s.userID, s.credentialID).
					Return(&repoCred, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedLogin, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			result, err := setup.service.Get(ctx, setup.userID, setup.encryptionKeyStr, setup.credentialID)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, result)
			} else if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assertCredentialFieldsWithID(t, result, setup.userID, setup.credentialID)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantCount int
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "successful list with multiple credentials",
			setupMock: func(s *testSetup) {
				cred1 := interfaces.RepositoryCredential{
					ID:                uuid.New().String(),
					UserID:            s.userID,
					Name:              "cred1",
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}
				cred2 := interfaces.RepositoryCredential{
					ID:                uuid.New().String(),
					UserID:            s.userID,
					Name:              "cred2",
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					ListCredentials(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryCredential{cred1, cred2}, nil).
					Times(1)

				// Expect decryption for each credential
				s.expectDecryptCredential()
				s.expectDecryptCredential()
			},
			wantCount: 2,
		},
		{
			name: "successful list with empty result",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					ListCredentials(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryCredential{}, nil).
					Times(1)
			},
			wantCount: 0,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					ListCredentials(gomock.Any(), s.userID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "decryption fails for one credential",
			setupMock: func(s *testSetup) {
				cred := interfaces.RepositoryCredential{
					ID:                uuid.New().String(),
					UserID:            s.userID,
					Name:              "cred1",
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					ListCredentials(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryCredential{cred}, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedLogin, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			result, err := setup.service.List(ctx, setup.userID, setup.encryptionKeyStr)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, result)
			} else if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "successful update",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				repoCred := interfaces.RepositoryCredential{
					ID:                s.credentialID,
					UserID:            s.userID,
					Name:              testName,
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					UpdateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCred, nil).
					Times(1)

				s.expectDecryptCredential()
			},
		},
		{
			name: "credential not found",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				s.mockRepo.EXPECT().
					UpdateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, ErrNotFound).
					Times(1)
			},
			wantErrIs: ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				s.mockRepo.EXPECT().
					UpdateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "repository returns nil",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				s.mockRepo.EXPECT().
					UpdateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "encryption fails",
			setupMock: func(s *testSetup) {
				s.mockCrypto.EXPECT().
					Encrypt(testLogin, s.encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "decryption fails after update",
			setupMock: func(s *testSetup) {
				s.expectEncryptCredential()

				repoCred := interfaces.RepositoryCredential{
					ID:                s.credentialID,
					UserID:            s.userID,
					Name:              testName,
					EncryptedLogin:    s.encryptedLogin,
					EncryptedPassword: s.encryptedPass,
					EncryptedNotes:    s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					UpdateCredential(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCred, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedLogin, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			cred := newTestCredential()
			cred.id = setup.credentialID
			cred.userID = setup.userID

			result, err := setup.service.Update(ctx, setup.userID, setup.encryptionKeyStr, *cred)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, result)
			} else if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assertCredentialFieldsWithID(t, result, setup.userID, setup.credentialID)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "successful deletion",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					DeleteCredential(gomock.Any(), s.userID, s.credentialID).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "credential not found",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					DeleteCredential(gomock.Any(), s.userID, s.credentialID).
					Return(ErrNotFound).
					Times(1)
			},
			wantErrIs: ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					DeleteCredential(gomock.Any(), s.userID, s.credentialID).
					Return(errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			err := setup.service.Delete(ctx, setup.userID, setup.credentialID)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			} else if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
