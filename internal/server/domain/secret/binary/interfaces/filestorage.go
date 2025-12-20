package interfaces

import (
	"context"
	"io"
)

type FileStorage interface {
	Upload(ctx context.Context, userID, fileID string, reader io.Reader) (filepath string, err error)
	Delete(ctx context.Context, userID, fileID string) error
}
