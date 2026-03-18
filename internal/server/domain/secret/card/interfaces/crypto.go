package interfaces

//go:generate mockgen -source=crypto.go -destination=mocks/mock_crypto.go -package=mocks

type CryptoService interface {
	Encrypt(plaintext string, key []byte) ([]byte, error)
	Decrypt(ciphertext []byte, key []byte) (string, error)
}
