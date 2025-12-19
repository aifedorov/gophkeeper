package interfaces

import "time"

type RepositoryFile struct {
	ID         string
	UserID     string
	Name       string
	Path       string
	Size       int64
	MimeType   string
	UploadedAt time.Time
}

type FileMetadata struct {
	Name     string
	Size     int64
	MimeType string
}
