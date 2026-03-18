package crypto

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

type encryptReader struct {
	source io.Reader
	gcm    cipher.AEAD
	nonce  []byte
	first  bool
	outBuf []byte
	outPos int
	eof    bool
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
		if len(p) < len(r.nonce) {
			return 0, io.ErrShortBuffer
		}
		return copy(p, r.nonce), nil
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

	plaintext := make([]byte, chunkSize)
	totalRead := 0
	for totalRead < chunkSize {
		readN, readErr := r.source.Read(plaintext[totalRead:])
		totalRead += readN
		if errors.Is(readErr, io.EOF) {
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

	ciphertext := r.gcm.Seal(nil, r.nonce, plaintext[:totalRead], nil)
	incrementNonce(r.nonce)

	n = copy(p, ciphertext)
	if n < len(ciphertext) {
		r.outBuf = ciphertext
		r.outPos = n
	}

	if r.eof && n >= len(ciphertext) {
		return n, io.EOF
	}
	return n, nil
}
