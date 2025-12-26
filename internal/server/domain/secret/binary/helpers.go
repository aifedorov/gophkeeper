package binary

import (
	"encoding/base64"
	"fmt"
	"io"
	"sync"

	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/crypto"
)

const keyAES256Len = 32 // AES-256 use 32 bytes key

func (s *service) wrapToEncrypt(encryptionKey string, reader io.Reader) (res io.Reader, key []byte, err error) {
	key, err = base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	if len(key) != keyAES256Len {
		return nil, nil, fmt.Errorf("invalid encryption key length: %d", len(key))
	}

	res, err = crypto.NewEncryptReader(reader, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt file: %w", err)
	}
	return res, key, nil
}

// acquireLock acquires an exclusive write lock for the given file ID.
// Use this for operations that modify the file (Update, Delete).
// Creates a new RWMutex if one doesn't exist for this file.
// The caller MUST call releaseLock() when done, preferably using defer.
func (s *service) acquireLock(fileID string) *sync.RWMutex {
	value, _ := s.locks.LoadOrStore(fileID, &sync.RWMutex{})
	mu := value.(*sync.RWMutex)
	mu.Lock()
	return mu
}

// acquireReadLock acquires a shared read lock for the given file ID.
// Use this for operations that only read the file (Download).
// Multiple readers can hold the lock simultaneously.
// The caller MUST call releaseReadLock() when done, preferably using defer.
func (s *service) acquireReadLock(fileID string) *sync.RWMutex {
	value, _ := s.locks.LoadOrStore(fileID, &sync.RWMutex{})
	mu := value.(*sync.RWMutex)
	mu.RLock()
	return mu
}

// releaseLock releases an exclusive write lock.
// Should always be called with defer after acquireLock().
func (s *service) releaseLock(mu *sync.RWMutex) {
	mu.Unlock()
}

// releaseReadLock releases a shared read lock.
// Should always be called with defer after acquireReadLock().
func (s *service) releaseReadLock(mu *sync.RWMutex) {
	mu.RUnlock()
}
