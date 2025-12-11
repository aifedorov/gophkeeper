package crypto

import (
	"crypto/rand"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	argonKeyLen  = 32 // AES-256
	saltLen      = 32
)

type cryptoService struct{}

func NewService() interfaces.CryptoService {
	return &cryptoService{}
}

func (s *cryptoService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func (s *cryptoService) DeriveEncryptionKey(password, salt string) []byte {
	return argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
}

func (s *cryptoService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (s *cryptoService) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
