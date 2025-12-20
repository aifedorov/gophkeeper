package binary

import (
	"context"
	"fmt"
	"os"

	authinterfaces "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/pkg/filestorage"
)

type Service interface {
	Upload(ctx context.Context, filePath string, notes string) error
	List(ctx context.Context) ([]File, error)
	Download(ctx context.Context, id string) (filepath string, error error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	client          Client
	store           *filestorage.FileStorage
	sessionProvider authinterfaces.SessionProvider
}

func NewService(client Client, store *filestorage.FileStorage, sessionProvider authinterfaces.SessionProvider) Service {
	return &service{
		client:          client,
		store:           store,
		sessionProvider: sessionProvider,
	}
}

func (s *service) Upload(ctx context.Context, filePath string, notes string) error {
	// #nosec G304
	f, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	fileInfo, err := NewFileInfo(f, notes)
	if err != nil {
		return fmt.Errorf("failed to create file info: %w", err)
	}

	return s.client.Upload(ctx, fileInfo, f)
}

func (s *service) List(ctx context.Context) ([]File, error) {
	return s.client.List(ctx)
}

func (s *service) Download(ctx context.Context, id string) (filepath string, error error) {
	session, err := s.sessionProvider.GetSession(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	reader, meta, err := s.client.Download(ctx, id)
	if err != nil || meta == nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	return s.store.Upload(ctx, session.GetUserID(), meta.Name(), reader)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, id)
}
