package interfaces

import "time"

type RepositoryFile struct {
	ID             string
	UserID         string
	Name           string
	EncryptedPath  []byte
	EncryptedSize  []byte
	EncryptedNotes []byte
	UploadedAt     time.Time
}

type FileMetadata struct {
	Name  string
	Size  int64
	Notes string
}
