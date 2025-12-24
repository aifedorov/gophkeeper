package auth

import (
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
)

func TestService_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*testSetup)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful login",
			setup: func(s *testSetup) {
				s.expectLoginSuccess()
			},
			wantErr: false,
		},
		{
			name: "login fails - client error",
			setup: func(s *testSetup) {
				s.expectLoginClientError(ErrInvalidCredentials)
			},
			wantErr:     true,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name: "login succeeds but save fails",
			setup: func(s *testSetup) {
				s.expectLoginSaveError(errors.New("storage error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.Login(s.ctx, s.testCreds)

			assertError(t, err, tt.wantErr, tt.expectedErr)
		})
	}
}

func TestService_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*testSetup)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful registration",
			setup: func(s *testSetup) {
				s.expectRegisterSuccess()
			},
			wantErr: false,
		},
		{
			name: "registration fails - auth already exists",
			setup: func(s *testSetup) {
				s.expectRegisterClientError(ErrUserAlreadyExists)
			},
			wantErr:     true,
			expectedErr: ErrUserAlreadyExists,
		},
		{
			name: "registration succeeds but save fails",
			setup: func(s *testSetup) {
				s.expectRegisterSaveError(errors.New("storage error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.Register(s.ctx, s.testCreds)

			assertError(t, err, tt.wantErr, tt.expectedErr)
		})
	}
}

func TestService_Logout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*testSetup)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful logout",
			setup: func(s *testSetup) {
				s.expectLogoutSuccess()
			},
			wantErr: false,
		},
		{
			name: "logout fails - delete error",
			setup: func(s *testSetup) {
				s.expectLogoutError(errors.New("delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.Logout(s.ctx)

			assertError(t, err, tt.wantErr, tt.expectedErr)
		})
	}
}

func TestService_GetCurrentSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*testSetup)
		wantSession shared.Session
		wantErr     bool
		expectedErr error
	}{
		{
			name: "get session successfully",
			setup: func(s *testSetup) {
				s.expectGetSessionSuccess()
			},
			wantSession: shared.NewSession(testToken, testToken, testUserID, testLogin),
			wantErr:     false,
		},
		{
			name: "session not found",
			setup: func(s *testSetup) {
				s.expectGetSessionError(ErrSessionNotFound)
			},
			wantErr:     true,
			expectedErr: ErrSessionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			session, err := s.service.GetCurrentSession()

			assertError(t, err, tt.wantErr, tt.expectedErr)
			if !tt.wantErr {
				assertSessionEqual(t, session, tt.wantSession)
			}
		})
	}
}
