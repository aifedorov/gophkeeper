package binary

import (
	"fmt"
)

type FileMeta struct {
	name  string
	size  int64
	notes string
}

func NewFileMeta(name string, size int64, notes string) (*FileMeta, error) {
	if name == "" {
		return nil, fmt.Errorf("file name is required")
	}
	if size == 0 {
		return nil, fmt.Errorf("file size can't be zero")
	}

	return &FileMeta{
		name:  name,
		size:  size,
		notes: notes,
	}, nil
}

func (f *FileMeta) Name() string {
	return f.name
}

func (f *FileMeta) Size() int64 {
	return f.size
}

func (f *FileMeta) Notes() string {
	return f.notes
}
