package crypto

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

const chunkSize = 64 * 1024

type encryptReader struct {
	source io.Reader
	gcm    cipher.AEAD
	nonce  []byte
	first  bool
}

func NewEncryptReader(reader io.Reader, key []byte) (io.Reader, error) {
	gcm, nonce, err := NewEncrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create encrypt data: %w", err)
	}

	return &encryptReader{
		source: reader,
		nonce:  nonce,
		gcm:    gcm,
		first:  true,
	}, nil
}

func (r *encryptReader) Read(p []byte) (n int, err error) {
	if r.first {
		r.first = false
		n := copy(p, r.nonce)
		if n < len(r.nonce) {
			return 0, io.ErrShortBuffer
		}
		return n, nil
	}

	maxSize := len(p) - r.gcm.Overhead()
	if maxSize < 0 {
		return 0, io.ErrShortBuffer
	}
	if maxSize > chunkSize {
		maxSize = chunkSize
	}

	plaintext := make([]byte, maxSize)
	n, err = r.source.Read(plaintext)
	if errors.Is(err, io.EOF) || n == 0 {
		return n, io.EOF
	}
	if err != nil {
		return n, fmt.Errorf("failed to read from source: %w", err)
	}

	ciphertext := r.gcm.Seal(nil, r.nonce, plaintext[:n], nil)
	copy(p, ciphertext)
	r.incrementNonce(r.nonce)

	return len(ciphertext), nil
}

func (r *encryptReader) incrementNonce(nonce []byte) {
	for i := len(nonce) - 1; i >= 0; i-- {
		nonce[i]++
		if nonce[i] != 0 {
			break
		}
	}
}
