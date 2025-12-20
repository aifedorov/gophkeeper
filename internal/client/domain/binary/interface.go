package binary

import (
	"context"
	"io"
)

type Client interface {
	Upload(ctx context.Context, fileInfo *FileMeta, reader io.Reader) error
	List(ctx context.Context) ([]File, error)
}
