// Package interfaces provides data models for binary file management.
package interfaces

import "time"

// RepositoryFile represents a file as stored in the repository.
// All sensitive fields (path, size, notes) are encrypted before storage.
type RepositoryFile struct {
	ID             string    // Unique file identifier
	UserID         string    // ID of the user who owns this file
	Name           string    // File name (stored in plain text for indexing)
	EncryptedPath  []byte    // Encrypted file storage path
	EncryptedSize  []byte    // Encrypted file size (as string)
	EncryptedNotes []byte    // Encrypted file notes/metadata
	Version        int64     // Version number
	UpdatedAt      time.Time // Timestamp when the file was updated
}

// FileMetadata represents file metadata used for upload and download operations.
type FileMetadata struct {
	ID      string // Unique file identifier
	Name    string // File name
	Size    int64  // File size in bytes
	Notes   string // Optional notes/metadata
	Version int64
}
