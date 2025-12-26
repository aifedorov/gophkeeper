package interfaces

type CacheStorage interface {
	SetFileVersion(id string, version int64) error
	GetFileVersion(id string) (int64, error)
	DeleteFileVersion(id string) error
}
