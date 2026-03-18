package shared

import (
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
)

func ParseErrorForCLI(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, shared.ErrUnauthenticated) {
		return fmt.Errorf("use register or login command to authenticate")
	}
	if errors.Is(err, shared.ErrAlreadyExists) {
		return fmt.Errorf("credential with such name already exists")
	}
	if errors.Is(err, shared.ErrNotFound) {
		return fmt.Errorf("credential with such ID not found")
	}
	if errors.Is(err, shared.ErrVersionConflict) {
		return fmt.Errorf("version conflict fetch data and try again")
	}
	return fmt.Errorf("unknown error %w", err)
}
