// Package binary provides binary file management services for the GophKeeper client.
//
// This package implements the client-side logic for managing binary file storage.
// It handles file uploads, downloads, listing, and deletion, communicating with the server via gRPC.
package binary

import (
	"context"
	"fmt"

	authinterfaces "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/domain/binary/interfaces"
)

// Service defines the interface for client-side binary file management operations.
type Service interface {
	// Upload sends a request to upload a binary file to the server.
	// It opens the file at the specified path, creates file metadata, and streams the file to the server.
	// Returns an error if the file doesn't exist, can't be opened, or if the upload fails.
	Upload(ctx context.Context, filePath string, notes string) error
	// List retrieves all binary files for the authenticated user from the server.
	// Returns an empty slice if the user has no files.
	List(ctx context.Context) ([]File, error)
	// Download retrieves a binary file by ID from the server and saves it to local storage.
	// It gets the user session to determine the storage location, downloads the file,
	// and saves it using the local file storage. Returns the path where the file was saved.
	// Returns an error if the session is not found, if the download fails, or if file save fails.
	Download(ctx context.Context, id string) (filepath string, error error)
	// Update sends a request to update an existing binary file on the server.
	// It opens the file at the specified path, creates file metadata with the given ID, and streams to the server.
	// Returns an error if the file doesn't exist, can't be opened, or if the update fails.
	Update(ctx context.Context, id string, filePath string, notes string) error
	// Delete sends a request to delete a binary file by ID from the server.
	// Returns an error if the file doesn't exist or if the deletion fails.
	Delete(ctx context.Context, id string) error
}

// service implements the Service interface for client-side binary file management.
type service struct {
	client          Client
	store           interfaces.Storage
	cache           interfaces.CacheStorage
	sessionProvider authinterfaces.SessionProvider
}

// NewService creates a new instance of the binary file service with the provided dependencies.
// It initializes the service with a gRPC client, local file storage, and session provider.
func NewService(client Client, store interfaces.Storage, cache interfaces.CacheStorage, sessionProvider authinterfaces.SessionProvider) Service {
	return &service{
		client:          client,
		store:           store,
		cache:           cache,
		sessionProvider: sessionProvider,
	}
}

// Upload sends a request to upload a binary file to the server.
// It opens the file at the specified path, creates file metadata, and streams the file to the server.
// Returns an error if the file doesn't exist, can't be opened, or if the upload fails.
func (s *service) Upload(ctx context.Context, filePath string, notes string) error {
	f, err := s.store.OpenFile(ctx, filePath)
	if err != nil {
		return err
	}

	fileInfo, err := NewFileInfo(f, notes)
	if err != nil {
		return fmt.Errorf("failed to create file info: %w", err)
	}

	id, version, err := s.client.Upload(ctx, fileInfo, f)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	err = s.cache.SetFileVersion(id, version)
	if err != nil {
		return fmt.Errorf("failed to save file version to cache: %w", err)
	}

	return nil
}

// List retrieves all binary files for the authenticated user from the server.
// Returns an empty slice if the user has no files.
func (s *service) List(ctx context.Context) ([]File, error) {
	files, err := s.client.List(ctx)
	if err != nil {
		return []File{}, fmt.Errorf("failed to list files: %w", err)
	}

	for _, file := range files {
		err := s.cache.SetFileVersion(file.ID(), file.Version())
		if err != nil {
			return []File{}, fmt.Errorf("failed to save file version to cache: %w", err)
		}
	}

	return files, nil
}

// Download retrieves a binary file by ID from the server and saves it to local storage.
// It gets the user session to determine the storage location, downloads the file,
// and saves it using the local file storage. Returns the path where the file was saved.
// Returns an error if the session is not found, if the download fails, or if file save fails.
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

	return s.store.Upload(ctx, session.GetLogin(), meta.Name(), reader)
}

// Update sends a request to update an existing binary file on the server.
// It opens the file at the specified path, creates file metadata with the given ID, and streams to the server.
// Returns an error if the file doesn't exist, can't be opened, or if the update fails.
func (s *service) Update(ctx context.Context, id string, filePath string, notes string) error {
	f, err := s.store.OpenFile(ctx, filePath)
	if err != nil {
		return err
	}

	currentVersion, err := s.cache.GetFileVersion(id)
	if err != nil {
		return fmt.Errorf("failed to get version from cache (try running 'list' first): %w", err)
	}

	fileInfo, err := NewUpdateFileInfoWithVersion(id, f, notes, currentVersion)
	if err != nil {
		return fmt.Errorf("failed to create update file info: %w", err)
	}

	newVersion, err := s.client.Update(ctx, fileInfo, f)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	err = s.cache.SetFileVersion(id, newVersion)
	if err != nil {
		return fmt.Errorf("failed to save file version to cache: %w", err)
	}

	return nil
}

// Delete sends a request to delete a binary file by ID from the server.
// Returns an error if the file doesn't exist or if the deletion fails.
func (s *service) Delete(ctx context.Context, id string) error {
	err := s.client.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	err = s.cache.DeleteFileVersion(id)
	if err != nil {
		return fmt.Errorf("failed to delete file version from cache: %w", err)
	}

	return nil
}
