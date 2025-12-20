package binary

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
)

func MetadataToDomain(metadata interfaces.FileMetadata) (*interfaces.File, error) {
	file, err := interfaces.NewFile(
		uuid.NewString(),
		metadata.Name,
		metadata.Size,
		metadata.Notes,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	return file, nil
}

func FileToRepository(crypto interfaces.CryptoService, key []byte, file *interfaces.File, path string) (interfaces.RepositoryFile, error) {
	encryptedPath, err := crypto.Encrypt(path, key)
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
		UploadedAt:     file.GetUploadedAt(),
	}, nil
}

func FileToDomain(crypto interfaces.CryptoService, key []byte, file interfaces.RepositoryFile) (*interfaces.File, error) {
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

	return interfaces.NewFile(
		file.ID,
		file.Name,
		int64(size),
		notes,
		file.UploadedAt,
	)
}
