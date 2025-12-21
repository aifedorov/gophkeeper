package card

import (
	"context"
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces"
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
				s.expectEncryptCard()

				repoCard := interfaces.RepositoryCard{
					ID:                      s.cardID,
					UserID:                  s.userID,
					Name:                    testName,
					EncryptedNumber:         s.encryptedNumber,
					EncryptedExpiredDate:    s.encryptedExpiredDate,
					EncryptedCardHolderName: s.encryptedCardHolderName,
					EncryptedCvv:            s.encryptedCvv,
					EncryptedNotes:          s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					CreateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCard, nil).
					Times(1)

				s.expectDecryptCard()
			},
		},
		{
			name: "name already exists",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				s.mockRepo.EXPECT().
					CreateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, ErrNameExists).
					Times(1)
			},
			wantErrIs: ErrNameExists,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				s.mockRepo.EXPECT().
					CreateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "repository returns nil",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				s.mockRepo.EXPECT().
					CreateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "encryption fails on number",
			setupMock: func(s *testSetup) {
				s.mockCrypto.EXPECT().
					Encrypt(testNumber, s.encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "decryption fails after creation",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				repoCard := interfaces.RepositoryCard{
					ID:                      s.cardID,
					UserID:                  s.userID,
					Name:                    testName,
					EncryptedNumber:         s.encryptedNumber,
					EncryptedExpiredDate:    s.encryptedExpiredDate,
					EncryptedCardHolderName: s.encryptedCardHolderName,
					EncryptedCvv:            s.encryptedCvv,
					EncryptedNotes:          s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					CreateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCard, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNumber, s.encryptionKey).
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
			card := newTestCard()
			result, err := setup.service.Create(ctx, setup.userID, setup.encryptionKeyStr, *card)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, result)
			} else if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assertCardFieldsWithID(t, result, setup.userID, setup.cardID)
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
			name: "successful list with multiple cards",
			setupMock: func(s *testSetup) {
				card1 := interfaces.RepositoryCard{
					ID:                      uuid.New().String(),
					UserID:                  s.userID,
					Name:                    "card1",
					EncryptedNumber:         s.encryptedNumber,
					EncryptedExpiredDate:    s.encryptedExpiredDate,
					EncryptedCardHolderName: s.encryptedCardHolderName,
					EncryptedCvv:            s.encryptedCvv,
					EncryptedNotes:          s.encryptedNotes,
				}
				card2 := interfaces.RepositoryCard{
					ID:                      uuid.New().String(),
					UserID:                  s.userID,
					Name:                    "card2",
					EncryptedNumber:         s.encryptedNumber,
					EncryptedExpiredDate:    s.encryptedExpiredDate,
					EncryptedCardHolderName: s.encryptedCardHolderName,
					EncryptedCvv:            s.encryptedCvv,
					EncryptedNotes:          s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					ListCards(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryCard{card1, card2}, nil).
					Times(1)

				s.expectDecryptCard()
				s.expectDecryptCard()
			},
			wantCount: 2,
		},
		{
			name: "successful list with empty result",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					ListCards(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryCard{}, nil).
					Times(1)
			},
			wantCount: 0,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					ListCards(gomock.Any(), s.userID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "decryption fails for one card",
			setupMock: func(s *testSetup) {
				card := interfaces.RepositoryCard{
					ID:                      uuid.New().String(),
					UserID:                  s.userID,
					Name:                    "card1",
					EncryptedNumber:         s.encryptedNumber,
					EncryptedExpiredDate:    s.encryptedExpiredDate,
					EncryptedCardHolderName: s.encryptedCardHolderName,
					EncryptedCvv:            s.encryptedCvv,
					EncryptedNotes:          s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					ListCards(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryCard{card}, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNumber, s.encryptionKey).
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
				s.expectEncryptCard()

				repoCard := interfaces.RepositoryCard{
					ID:                      s.cardID,
					UserID:                  s.userID,
					Name:                    testName,
					EncryptedNumber:         s.encryptedNumber,
					EncryptedExpiredDate:    s.encryptedExpiredDate,
					EncryptedCardHolderName: s.encryptedCardHolderName,
					EncryptedCvv:            s.encryptedCvv,
					EncryptedNotes:          s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					UpdateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCard, nil).
					Times(1)

				s.expectDecryptCard()
			},
		},
		{
			name: "card not found",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				s.mockRepo.EXPECT().
					UpdateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, ErrNotFound).
					Times(1)
			},
			wantErrIs: ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				s.mockRepo.EXPECT().
					UpdateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "repository returns nil",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				s.mockRepo.EXPECT().
					UpdateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "encryption fails",
			setupMock: func(s *testSetup) {
				s.mockCrypto.EXPECT().
					Encrypt(testNumber, s.encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "decryption fails after update",
			setupMock: func(s *testSetup) {
				s.expectEncryptCard()

				repoCard := interfaces.RepositoryCard{
					ID:                      s.cardID,
					UserID:                  s.userID,
					Name:                    testName,
					EncryptedNumber:         s.encryptedNumber,
					EncryptedExpiredDate:    s.encryptedExpiredDate,
					EncryptedCardHolderName: s.encryptedCardHolderName,
					EncryptedCvv:            s.encryptedCvv,
					EncryptedNotes:          s.encryptedNotes,
				}

				s.mockRepo.EXPECT().
					UpdateCard(gomock.Any(), s.userID, gomock.Any()).
					Return(&repoCard, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNumber, s.encryptionKey).
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
			card := newTestCard()
			card.id = setup.cardID
			card.userID = setup.userID

			result, err := setup.service.Update(ctx, setup.userID, setup.encryptionKeyStr, *card)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, result)
			} else if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assertCardFieldsWithID(t, result, setup.userID, setup.cardID)
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
					DeleteCard(gomock.Any(), s.userID, s.cardID).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "card not found",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					DeleteCard(gomock.Any(), s.userID, s.cardID).
					Return(ErrNotFound).
					Times(1)
			},
			wantErrIs: ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					DeleteCard(gomock.Any(), s.userID, s.cardID).
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
			err := setup.service.Delete(ctx, setup.userID, setup.cardID)

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
