// Package text provides text note management services for the GophKeeper client.
//
// This package implements the client-side logic for managing text notes.
// It wraps the binary file service to provide a convenient interface for text operations,
// supporting both inline content and file-based input.
package text

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/aifedorov/gophkeeper/pkg/filestorage"
)

const (
	// MaxViewSize is the maximum file size (in bytes) that can be viewed inline.
	// Files larger than this will require explicit download.
	MaxViewSize = 100 * 1024 // 100KB
)

// Service defines the interface for client-side text note management operations.
type Service interface {
	// CreateFromContent creates a new text note from inline content.
	// It creates a temporary file, uploads it, and cleans up.
	CreateFromContent(ctx context.Context, content, title, notes string) error
	// CreateFromFile creates a new text note from an existing file.
	CreateFromFile(ctx context.Context, filePath, notes string) error
	// View retrieves and displays a text note by ID.
	// If the file is larger than MaxViewSize, it returns an error suggesting download.
	View(ctx context.Context, id string) (content string, err error)
	// List retrieves all text notes for the authenticated user.
	List(ctx context.Context) ([]binary.File, error)
	// UpdateFromContent updates an existing text note with new inline content.
	UpdateFromContent(ctx context.Context, id, content, title, notes string) error
	// UpdateFromFile updates an existing text note from a file.
	UpdateFromFile(ctx context.Context, id, filePath, notes string) error
	// Download retrieves a text note by ID and saves it to local storage.
	Download(ctx context.Context, id string) (filepath string, error error)
	// Delete removes a text note by ID from the server.
	Delete(ctx context.Context, id string) error
}

// service implements the Service interface for text note management.
type service struct {
	binarySrv binary.Service
	store     filestorage.Storage
}

// NewService creates a new instance of the text service.
func NewService(binarySrv binary.Service, store filestorage.Storage) Service {
	return &service{
		binarySrv: binarySrv,
		store:     store,
	}
}

// CreateFromContent creates a new text note from inline content.
// It creates a temporary file with the content, uploads it, and cleans up the temp file.
func (s *service) CreateFromContent(ctx context.Context, content, title, notes string) error {
	reader := strings.NewReader(content)
	path, err := s.store.Upload(ctx, ".tmp", fmt.Sprintf("%s.txt", title), reader)
	defer func() {
		_ = os.Remove(path)
	}()

	if err != nil {
		return fmt.Errorf("failed to create file from content: %w", err)
	}
	return s.binarySrv.Upload(ctx, path, notes)
}

// CreateFromFile creates a new text note from an existing file.
// This is a simple delegation to the binary service.
func (s *service) CreateFromFile(ctx context.Context, filePath, notes string) error {
	return s.binarySrv.Upload(ctx, filePath, notes)
}

// View retrieves and displays a text note by ID.
// If the file is larger than MaxViewSize, it returns an error.
func (s *service) View(ctx context.Context, id string) (string, error) {
	path, err := s.binarySrv.Download(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	content, err := s.store.ReadContent(ctx, path, MaxViewSize)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}
	return content, nil
}

// List retrieves all text notes for the authenticated user.
func (s *service) List(ctx context.Context) ([]binary.File, error) {
	return s.binarySrv.List(ctx)
}

// UpdateFromContent updates an existing text note with new inline content.
func (s *service) UpdateFromContent(ctx context.Context, id, content, title, notes string) error {
	reader := strings.NewReader(content)
	path, err := s.store.Upload(ctx, ".tmp", fmt.Sprintf("%s.txt", title), reader)
	defer func() {
		_ = os.Remove(path)
	}()
	if err != nil {
		return fmt.Errorf("failed to create file from content: %w", err)
	}
	return s.binarySrv.Update(ctx, id, path, notes)
}

// UpdateFromFile updates an existing text note from a file.
func (s *service) UpdateFromFile(ctx context.Context, id, filePath, notes string) error {
	return s.binarySrv.Update(ctx, id, filePath, notes)
}

// Download retrieves a text note by ID and saves it to local storage.
func (s *service) Download(ctx context.Context, id string) (string, error) {
	return s.binarySrv.Download(ctx, id)
}

// Delete removes a text note by ID from the server.
func (s *service) Delete(ctx context.Context, id string) error {
	return s.binarySrv.Delete(ctx, id)
}
