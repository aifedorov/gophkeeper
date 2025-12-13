package interfaces

//go:generate mockgen -source=crypto.go -destination=mocks/mock_crypto.go -package=mocks

// CryptoService defines the interface for cryptographic operations.
type CryptoService interface {
	// Encrypt encrypts a plaintext string using a key.
	// Returns the encrypted ciphertext or an error if encryption fails.
	Encrypt(plaintext string, key []byte) ([]byte, error)
	// Decrypt decrypts a ciphertext using a key.
	// Returns the decrypted plaintext or an error if decryption fails.
	Decrypt(ciphertext []byte, key []byte) (string, error)
}
