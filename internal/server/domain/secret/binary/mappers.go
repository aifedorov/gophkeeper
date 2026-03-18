// Package binary provides mappers for converting between domain and repository representations.
package binary

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
)

// MetadataToFile converts file metadata to a domain File entity.
// It generates a new UUID for the file if ID is empty, otherwise uses the provided ID.
// Sets the current time as the upload timestamp.
// For new files (version 0), sets version to 1.
// Returns an error if the metadata is invalid (e.g., empty name or zero size).
func MetadataToFile(metadata interfaces.FileMetadata) (*interfaces.File, error) {
	id := metadata.ID
	if id == "" {
		id = uuid.NewString()
	}
	version := metadata.Version
	if version == 0 {
		version = 1 // Default version for new files
	}
	file, err := interfaces.NewFile(
		id,
		metadata.Name,
		metadata.Size,
		"",
		metadata.Notes,
		version,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	return file, nil
}

// FileToRepository converts a domain File entity to a repository file representation.
// It encrypts all sensitive fields (path, size, notes) using the provided encryption key.
// Returns an error if the file is nil or if encryption fails for any field.
func FileToRepository(crypto interfaces.CryptoService, key []byte, file *interfaces.File) (interfaces.RepositoryFile, error) {
	if file == nil {
		return interfaces.RepositoryFile{}, fmt.Errorf("file is nil")
	}

	encryptedPath, err := crypto.Encrypt(file.GetPath(), key)
	if err != nil {
		return interfaces.RepositoryFile{}, fmt.Errorf("failed to encrypt path: %w", err)
	}
	encryptedNotes, err := crypto.Encrypt(file.GetNotes(), key)
	if err != nil {
		return interfaces.RepositoryFile{}, fmt.Errorf("failed to encrypt notes: %w", err)
	}
	encryptedSize, err := crypto.Encrypt(fmt.Sprintf("%d", file.GetSize()), key)
	if err != nil {
		return interfaces.RepositoryFile{}, fmt.Errorf("failed to encrypt size: %w", err)
	}

	return interfaces.RepositoryFile{
		ID:             file.GetID(),
		Name:           file.GetName(),
		EncryptedPath:  encryptedPath,
		EncryptedSize:  encryptedSize,
		EncryptedNotes: encryptedNotes,
		Version:        file.GetVersion(),
		UpdatedAt:      file.GetUploadedAt(),
	}, nil
}

// RepositoryToDomain converts a repository file representation to a domain File entity.
// It decrypts all encrypted fields (path, size, notes) using the provided encryption key.
// Returns an error if decryption fails for any field or if size conversion fails.
func RepositoryToDomain(crypto interfaces.CryptoService, key []byte, file *interfaces.RepositoryFile) (*interfaces.File, error) {
	if file == nil {
		return nil, fmt.Errorf("file is nil")
	}

	notes, err := crypto.Decrypt(file.EncryptedNotes, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt notes: %w", err)
	}
	sizeStr, err := crypto.Decrypt(file.EncryptedSize, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt size: %w", err)
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert size: %w", err)
	}
	path, err := crypto.Decrypt(file.EncryptedPath, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt path: %w", err)
	}

	return interfaces.NewFile(
		file.ID,
		file.Name,
		int64(size),
		path,
		notes,
		file.Version,
		file.UpdatedAt,
	)
}

// FileToMetadata converts a domain File entity to file metadata.
// Returns an error if the file is nil.
func FileToMetadata(file *interfaces.File) (interfaces.FileMetadata, error) {
	if file == nil {
		return interfaces.FileMetadata{}, fmt.Errorf("file is nil")
	}

	return interfaces.FileMetadata{
		Name:    file.GetName(),
		Size:    file.GetSize(),
		Notes:   file.GetNotes(),
		Version: file.GetVersion(),
	}, nil
}
