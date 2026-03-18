package credential

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testIDCred       = "test-id-123"
	testNameCred     = "test-credential"
	testLoginCred    = "testuser"
	testPasswordCred = "testpassword"
	testNotesCred    = "test notes"
	testVersionCred  = int64(1)
)

func TestNewCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		id         string
		credName   string
		login      string
		password   string
		metadata   string
		version    int64
		wantErr    error
		wantErrStr string
	}{
		{
			name:     "creates credential with all fields",
			id:       testIDCred,
			credName: testNameCred,
			login:    testLoginCred,
			password: testPasswordCred,
			metadata: testNotesCred,
			version:  testVersionCred,
		},
		{
			name:     "creates credential without metadata",
			id:       "id-2",
			credName: "cred-2",
			login:    "user2",
			password: "pass2",
			metadata: "",
			version:  2,
		},
		{
			name:     "empty id",
			id:       "",
			credName: testNameCred,
			login:    testLoginCred,
			password: testPasswordCred,
			metadata: testNotesCred,
			version:  testVersionCred,
			wantErr:  ErrIDRequired,
		},
		{
			name:     "empty name",
			id:       testIDCred,
			credName: "",
			login:    testLoginCred,
			password: testPasswordCred,
			metadata: testNotesCred,
			version:  testVersionCred,
			wantErr:  ErrNameRequired,
		},
		{
			name:     "empty login",
			id:       testIDCred,
			credName: testNameCred,
			login:    "",
			password: testPasswordCred,
			metadata: testNotesCred,
			version:  testVersionCred,
			wantErr:  ErrLoginRequired,
		},
		{
			name:     "empty password",
			id:       testIDCred,
			credName: testNameCred,
			login:    testLoginCred,
			password: "",
			metadata: testNotesCred,
			version:  testVersionCred,
			wantErr:  ErrPasswordRequired,
		},
		{
			name:       "zero version",
			id:         testIDCred,
			credName:   testNameCred,
			login:      testLoginCred,
			password:   testPasswordCred,
			metadata:   testNotesCred,
			version:    0,
			wantErrStr: "invalid credential version",
		},
		{
			name:       "negative version",
			id:         testIDCred,
			credName:   testNameCred,
			login:      testLoginCred,
			password:   testPasswordCred,
			metadata:   testNotesCred,
			version:    -1,
			wantErrStr: "invalid credential version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cred, err := NewCredential(tt.id, tt.credName, tt.login, tt.password, tt.metadata, tt.version)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, cred)
				return
			}
			if tt.wantErrStr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrStr)
				assert.Nil(t, cred)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cred)
			assert.Equal(t, tt.id, cred.GetID())
			assert.Equal(t, tt.credName, cred.GetName())
			assert.Equal(t, tt.login, cred.GetLogin())
			assert.Equal(t, tt.password, cred.GetPassword())
			assert.Equal(t, tt.metadata, cred.GetMetadata())
			assert.Equal(t, tt.version, cred.GetVersion())
		})
	}
}

func TestCredential_Getters(t *testing.T) {
	t.Parallel()

	cred, err := NewCredential(testIDCred, testNameCred, testLoginCred, testPasswordCred, testNotesCred, testVersionCred)
	require.NoError(t, err)

	assert.Equal(t, testIDCred, cred.GetID())
	assert.Equal(t, "", cred.GetUserID()) // userID is not set by NewCredential
	assert.Equal(t, testNameCred, cred.GetName())
	assert.Equal(t, testLoginCred, cred.GetLogin())
	assert.Equal(t, testPasswordCred, cred.GetPassword())
	assert.Equal(t, testNotesCred, cred.GetMetadata())
	assert.Equal(t, testVersionCred, cred.GetVersion())
}
