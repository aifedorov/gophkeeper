package binary

import (
	"time"
)

type File struct {
	id         string
	name       string
	size       int64
	notes      string
	uploadedAt time.Time
}

func NewFile(id, name string, size int64, notes string, uploadedAt time.Time) (*File, error) {
	return &File{id: id, name: name, size: size, notes: notes, uploadedAt: uploadedAt}, nil
}

func (f *File) ID() string {
	return f.id
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int64 {
	return f.size
}

func (f *File) Notes() string {
	return f.notes
}

func (f *File) UploadedAt() time.Time {
	return f.uploadedAt
}
