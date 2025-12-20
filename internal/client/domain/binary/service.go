package binary

import (
	"context"
	"fmt"
	"os"
)

type Service interface {
	Upload(ctx context.Context, filePath string) error
}

type service struct {
	client Client
}

func NewService(client Client) Service {
	return &service{
		client: client,
	}
}

func (s *service) Upload(ctx context.Context, filePath string) error {
	// #nosec G304
	f, err := os.Open(filePath)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("file not found: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	fileInfo, err := NewFileInfo(f)
	if err != nil {
		return fmt.Errorf("failed to create file info: %w", err)
	}

	return s.client.Upload(ctx, fileInfo, f)
}
