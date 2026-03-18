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
	outBuf []byte
	outPos int
	eof    bool
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

	if r.outPos < len(r.outBuf) {
		n = copy(p, r.outBuf[r.outPos:])
		r.outPos += n
		if r.outPos >= len(r.outBuf) {
			r.outBuf = nil
			r.outPos = 0
		}
		if r.eof && r.outPos >= len(r.outBuf) {
			return n, io.EOF
		}
		return n, nil
	}

	if r.eof {
		return 0, io.EOF
	}

	ciphertextSize := chunkSize + r.gcm.Overhead()
	ciphertext := make([]byte, ciphertextSize)

	totalRead := 0
	for totalRead < ciphertextSize {
		readN, readErr := r.source.Read(ciphertext[totalRead:])
		totalRead += readN
		if readErr == io.EOF {
			r.eof = true
			break
		}
		if readErr != nil {
			return 0, fmt.Errorf("failed to read from source: %w", readErr)
		}
	}

	if totalRead == 0 {
		return 0, io.EOF
	}

	plaintext, decErr := r.gcm.Open(nil, r.nonce, ciphertext[:totalRead], nil)
	if decErr != nil {
		return 0, fmt.Errorf("failed to decrypt: %w", decErr)
	}

	incrementNonce(r.nonce)

	n = copy(p, plaintext)
	if n < len(plaintext) {
		r.outBuf = plaintext
		r.outPos = n
	}

	if r.eof && n >= len(plaintext) {
		return n, io.EOF
	}
	return n, nil
}
