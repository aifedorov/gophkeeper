package binary

import (
	"context"
	"io"
)

type Client interface {
	Upload(ctx context.Context, fileInfo *FileInfo, reader io.Reader) error
	List(ctx context.Context) ([]File, error)
	Download(ctx context.Context, id string) (io.ReadCloser, *FileMeta, error)
	Delete(ctx context.Context, id string) error
}
