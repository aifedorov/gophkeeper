package binary

import (
	"context"
	"io"
)

type Client interface {
	Upload(ctx context.Context, fileInfo *FileInfo, reader io.Reader) error
}
