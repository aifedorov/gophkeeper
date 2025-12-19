package binary

import (
	"fmt"
	"time"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
)

func MetadataToDomain(metadata interfaces.FileMetadata) (*interfaces.File, error) {
	file, err := interfaces.NewFile(
		uuid.NewString(),
		metadata.Name,
		metadata.Size,
		metadata.MimeType,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	return file, nil
}

func DomainToRepository(file *interfaces.File, path string) interfaces.RepositoryFile {
	return interfaces.RepositoryFile{
		ID:         file.GetID(),
		Name:       file.GetName(),
		Path:       path,
		Size:       file.GetSize(),
		MimeType:   file.GetMimeType(),
		UploadedAt: file.GetUploadedAt(),
	}
}
