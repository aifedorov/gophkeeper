package interfaces

// CryptoService defines the interface for cryptographic operations.
type CryptoService interface {
	// GenerateSalt generates a random salt for password encryption.
	// Returns an error if salt generation fails.
	GenerateSalt() ([]byte, error)
	// DeriveEncryptionKey derives an encryption key from password and salt.
	// Uses Argon2ID key derivation function.
	DeriveEncryptionKey(password, salt string) []byte
	// HashPassword hashes a password using bcrypt.
	// Returns the hashed password as a string or an error if hashing fails.
	HashPassword(password string) (string, error)
	// CompareHashAndPassword compares a hashed password with a plaintext password.
	// Returns nil if they match, or an error if they don't match or comparison fails.
	CompareHashAndPassword(hashedPassword, password string) error
}
