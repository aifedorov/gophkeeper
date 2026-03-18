package card

import (
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
)

// propagateError wraps errors and preserves specific error types that need special handling.
// It ensures that authentication, conflict, and not-found errors are propagated without wrapping,
// while other errors are wrapped with the operation context.
func propagateError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, shared.ErrUnauthenticated) {
		return shared.ErrUnauthenticated
	}
	if errors.Is(err, shared.ErrAlreadyExists) {
		return shared.ErrAlreadyExists
	}
	if errors.Is(err, shared.ErrVersionConflict) {
		return shared.ErrVersionConflict
	}
	if errors.Is(err, shared.ErrNotFound) {
		return shared.ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
