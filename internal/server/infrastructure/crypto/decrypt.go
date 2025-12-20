package crypto

import (
	"crypto/cipher"
	"fmt"
	"io"
)

type decryptReader struct {
	source io.Reader
	gcm    cipher.AEAD
	nonce  []byte
	first  bool
}

func NewDecryptReader(reader io.Reader, key []byte) (io.Reader, error) {
	gcm, err := NewDecrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create decrypt data: %w", err)
	}

	return &decryptReader{
		source: reader,
		gcm:    gcm,
		nonce:  make([]byte, gcm.NonceSize()),
		first:  true,
	}, nil
}

func (r *decryptReader) Read(p []byte) (n int, err error) {
	if r.first {
		r.first = false
		n, err := io.ReadFull(r.source, r.nonce)
		if err != nil {
			return 0, fmt.Errorf("failed to read nonce: %w", err)
		}
		if n < len(r.nonce) {
			return 0, fmt.Errorf("incomplete nonce read")
		}
		return r.Read(p)
	}

	maxPlaintextSize := len(p) - r.gcm.Overhead()
	if maxPlaintextSize < 0 {
		return 0, io.ErrShortBuffer
	}
	if maxPlaintextSize > chunkSize {
		maxPlaintextSize = chunkSize
	}

	ciphertextSize := maxPlaintextSize + r.gcm.Overhead()
	ciphertext := make([]byte, ciphertextSize)

	n, err = r.source.Read(ciphertext)
	if err == io.EOF && n == 0 {
		return 0, io.EOF
	}
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("failed to read from source: %w", err)
	}

	plaintext, decErr := r.gcm.Open(nil, r.nonce, ciphertext[:n], nil)
	if decErr != nil {
		return 0, fmt.Errorf("failed to decrypt: %w", decErr)
	}

	incrementNonce(r.nonce)

	copied := copy(p, plaintext)
	if err == io.EOF {
		return copied, io.EOF
	}
	return copied, nil
}
