package interfaces

type CryptoService interface {
	Encrypt(plaintext string, key []byte) ([]byte, error)
	Decrypt(ciphertext []byte, key []byte) (string, error)
}
