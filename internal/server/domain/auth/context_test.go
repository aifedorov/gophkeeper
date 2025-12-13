package auth

import (
	"context"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestService_SetAndGetUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "set and get valid user ID",
			userID:  "test-user-123",
			wantErr: false,
		},
		{
			name:    "set and get UUID user ID",
			userID:  "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
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

			ctx := context.Background()
			ctxWithUserID := service.SetUserID(ctx, tt.userID)

			retrievedUserID, err := service.GetUserIDFromContext(ctxWithUserID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.userID, retrievedUserID)
			}
		})
	}
}

func TestService_GetUserIDFromContext_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "empty context",
			ctx:  context.Background(),
		},
		{
			name: "context with empty user ID",
			ctx:  context.WithValue(context.Background(), userIDKey, ""),
		},
		{
			name: "context with wrong type",
			ctx:  context.WithValue(context.Background(), userIDKey, 123),
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

			_, err := service.GetUserIDFromContext(tt.ctx)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to get user id from context")
		})
	}
}

func TestService_SetAndGetEncryptionKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		encryptionKey string
	}{
		{
			name:          "set and get base64 encryption key",
			encryptionKey: "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEh",
		},
		{
			name:          "set and get another encryption key",
			encryptionKey: "YW5vdGhlci1rZXktZm9yLXRlc3Rpbmc=",
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

			ctx := context.Background()
			ctxWithKey := service.SetEncryptionKeyEncoded(ctx, tt.encryptionKey)

			retrievedKey, err := service.GetEncryptionKeyFromContext(ctxWithKey)

			require.NoError(t, err)
			assert.Equal(t, tt.encryptionKey, retrievedKey)
		})
	}
}
