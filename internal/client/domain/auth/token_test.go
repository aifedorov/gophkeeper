package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewSessionProvider(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockSessionStore(ctrl)
	provider := NewSessionProvider(mockStore)

	require.NotNil(t, provider)
}

func TestSessionProvider_GetSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*MockSessionStore)
		wantSession shared.Session
		wantErr     error
	}{
		{
			name: "successful session retrieval",
			setup: func(m *MockSessionStore) {
				session := shared.NewSession("access-token", "encryption-key", "user-id", "login")
				m.EXPECT().Load().Return(session, nil).Times(1)
			},
			wantSession: shared.NewSession("access-token", "encryption-key", "user-id", "login"),
			wantErr:     nil,
		},
		{
			name: "session not found",
			setup: func(m *MockSessionStore) {
				m.EXPECT().Load().Return(shared.Session{}, errors.New("not found")).Times(1)
			},
			wantSession: shared.Session{},
			wantErr:     ErrSessionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := NewMockSessionStore(ctrl)
			tt.setup(mockStore)

			provider := NewSessionProvider(mockStore)
			session, err := provider.GetSession(context.Background())

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantSession, session)
		})
	}
}
