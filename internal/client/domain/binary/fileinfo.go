package binary

import (
	"fmt"
	"os"
)

type FileInfo struct {
	name  string
	size  int64
	notes string
}

func NewFileInfo(file *os.File, notes string) (*FileInfo, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	return &FileInfo{
		name:  stat.Name(),
		size:  stat.Size(),
		notes: notes,
	}, nil
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) Notes() string {
	return f.notes
}
