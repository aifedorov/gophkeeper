package repository

import (
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*testSetup)
		userID      string
		fileID      string
		wantErr     bool
		expectedErr error
		errMsg      string
	}{
		{
			name: "successful get",
			setup: func(s *testSetup) {
				expectedFile := s.newTestFile()
				s.expectGetSuccess(expectedFile)
			},
			wantErr: false,
		},
		{
			name: "file not found",
			setup: func(s *testSetup) {
				s.expectGetNotFound()
			},
			wantErr:     true,
			expectedErr: binary.ErrNotFound,
		},
		{
			name:    "invalid user UUID",
			setup:   func(s *testSetup) {},
			userID:  "invalid-uuid",
			wantErr: true,
			errMsg:  "failed to parse user id",
		},
		{
			name:    "invalid file UUID",
			setup:   func(s *testSetup) {},
			fileID:  "invalid-uuid",
			wantErr: true,
			errMsg:  "failed to parse file id",
		},
		{
			name: "database error",
			setup: func(s *testSetup) {
				s.expectGetError(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to get file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initRepoForGet()
			tt.setup(s)

			userID := s.userID.String()
			fileID := s.fileID.String()
			if tt.userID != "" {
				userID = tt.userID
			}
			if tt.fileID != "" {
				fileID = tt.fileID
			}

			result, err := s.repo.Get(s.ctx, userID, fileID)

			if tt.wantErr {
				assertError(t, err, true, tt.expectedErr, tt.errMsg)
				assert.Nil(t, result)
			} else {
				expectedFile := s.newTestFile()
				assertError(t, err, false, nil, "")
				assertFileEqual(t, *result, expectedFile)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*testSetup)
		userID      string
		fileID      string
		wantErr     bool
		expectedErr error
		errMsg      string
	}{
		{
			name: "successful update",
			setup: func(s *testSetup) {
				existingFile := s.newTestFile()
				newFile := s.newTestRepositoryFile()
				s.expectUpdateSuccess(existingFile, newFile)
			},
			wantErr: false,
		},
		{
			name: "file not found",
			setup: func(s *testSetup) {
				s.expectUpdateNotFound()
			},
			wantErr:     true,
			expectedErr: binary.ErrNotFound,
		},
		{
			name: "begin transaction error",
			setup: func(s *testSetup) {
				s.expectUpdateBeginError(errors.New("transaction error"))
			},
			wantErr: true,
			errMsg:  "failed to begin transaction",
		},
		{
			name: "update file error",
			setup: func(s *testSetup) {
				existingFile := s.newTestFile()
				s.expectUpdateFileError(existingFile, errors.New("update error"))
			},
			wantErr: true,
			errMsg:  "failed to update file",
		},
		{
			name:    "invalid user UUID",
			setup:   func(s *testSetup) {},
			userID:  "invalid-uuid",
			wantErr: true,
			errMsg:  "failed to parse user id",
		},
		{
			name:    "invalid file UUID",
			setup:   func(s *testSetup) {},
			fileID:  "invalid-uuid",
			wantErr: true,
			errMsg:  "failed to parse binary id",
		},
		{
			name: "commit error",
			setup: func(s *testSetup) {
				existingFile := s.newTestFile()
				s.expectUpdateCommitError(existingFile)
			},
			wantErr: true,
			errMsg:  "commit error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initRepoForUpdate()
			tt.setup(s)

			userID := s.userID.String()
			fileID := s.fileID.String()
			if tt.userID != "" {
				userID = tt.userID
			}
			if tt.fileID != "" {
				fileID = tt.fileID
			}

			newFile := s.newTestRepositoryFile()
			_, err := s.repo.Update(s.ctx, userID, fileID, newFile)

			assertError(t, err, tt.wantErr, tt.expectedErr, tt.errMsg)
		})
	}
}
