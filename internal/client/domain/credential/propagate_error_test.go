package credential

import (
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestPropagateError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation string
		err       error
		wantErr   error
		wantNil   bool
	}{
		{
			name:      "nil error returns nil",
			operation: "test",
			err:       nil,
			wantNil:   true,
		},
		{
			name:      "unauthenticated error is propagated",
			operation: "test",
			err:       shared.ErrUnauthenticated,
			wantErr:   shared.ErrUnauthenticated,
		},
		{
			name:      "wrapped unauthenticated error is propagated",
			operation: "test",
			err:       errors.Join(errors.New("wrapper"), shared.ErrUnauthenticated),
			wantErr:   shared.ErrUnauthenticated,
		},
		{
			name:      "already exists error is propagated",
			operation: "test",
			err:       shared.ErrAlreadyExists,
			wantErr:   shared.ErrAlreadyExists,
		},
		{
			name:      "version conflict error is propagated",
			operation: "test",
			err:       shared.ErrVersionConflict,
			wantErr:   shared.ErrVersionConflict,
		},
		{
			name:      "not found error is propagated",
			operation: "test",
			err:       shared.ErrNotFound,
			wantErr:   shared.ErrNotFound,
		},
		{
			name:      "other error is wrapped with operation",
			operation: "failed to do something",
			err:       errors.New("some error"),
			wantErr:   nil, // check that error contains operation and error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := propagateError(tt.operation, tt.err)

			if tt.wantNil {
				assert.Nil(t, result)
				return
			}

			if tt.wantErr != nil {
				assert.ErrorIs(t, result, tt.wantErr)
			} else {
				assert.Contains(t, result.Error(), tt.operation)
				assert.Contains(t, result.Error(), tt.err.Error())
			}
		})
	}
}
