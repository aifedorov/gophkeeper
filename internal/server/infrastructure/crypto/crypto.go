package crypto

import (
	"crypto/rand"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	argonTime     = 1
	argonMemory   = 64 * 1024 // 64 MB
	argonThreads  = 4
	argonKeyLen   = 32 // AES-256
	saltLen       = 32
	aes256KeySize = 32
)

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) GenerateSalt() ([]byte, error) {
	s.logger.Debug("crypto: generating salt")
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		s.logger.Error("crypto: failed to generate salt", zap.Error(err))
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	s.logger.Debug("crypto: salt generated successfully")
	return salt, nil
}

func (s *Service) DeriveEncryptionKey(password, salt string) []byte {
	s.logger.Debug("crypto: deriving encryption key using Argon2")
	key := argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
	s.logger.Debug("crypto: encryption key derived successfully")
	return key
}

func (s *Service) HashPassword(password string) (string, error) {
	s.logger.Debug("crypto: hashing password with bcrypt")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("crypto: failed to hash password", zap.Error(err))
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	s.logger.Debug("crypto: password hashed successfully")
	return string(hash), nil
}

func (s *Service) CompareHashAndPassword(hashedPassword, password string) error {
	s.logger.Debug("crypto: comparing password hash")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		s.logger.Debug("crypto: password hash comparison failed", zap.Error(err))
		return err
	}
	s.logger.Debug("crypto: password hash comparison successful")
	return nil
}

func (s *Service) Encrypt(plaintext string, key []byte) ([]byte, error) {
	s.logger.Debug("crypto: encrypting data", zap.Int("plaintext_len", len(plaintext)))

	if len(key) != aes256KeySize {
		s.logger.Error("crypto: invalid key size", zap.Int("key_size", len(key)))
		return nil, fmt.Errorf("key must be %d bytes for AES-256", aes256KeySize)
	}

	gcm, nonce, err := NewEncrypt(key)
	if err != nil {
		s.logger.Error("crypto: failed to create encrypt data", zap.Error(err))
		return nil, fmt.Errorf("failed to create encrypt data: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	s.logger.Debug("crypto: data encrypted successfully", zap.Int("ciphertext_len", len(ciphertext)))
	return ciphertext, nil
}

func (s *Service) Decrypt(ciphertext []byte, key []byte) (string, error) {
	s.logger.Debug("crypto: decrypting data", zap.Int("ciphertext_len", len(ciphertext)))

	if len(key) != aes256KeySize {
		s.logger.Error("crypto: invalid key size", zap.Int("key_size", len(key)))
		return "", fmt.Errorf("key must be %d bytes for AES-256", aes256KeySize)
	}

	gcm, err := NewDecrypt(key)
	if err != nil {
		s.logger.Error("crypto: failed to create decrypt data", zap.Error(err))
		return "", fmt.Errorf("failed to create decrypt data: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		s.logger.Error("crypto: ciphertext too short", zap.Int("ciphertext_len", len(ciphertext)), zap.Int("nonce_size", nonceSize))
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		s.logger.Error("crypto: failed to decrypt", zap.Error(err))
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	s.logger.Debug("crypto: data decrypted successfully", zap.Int("plaintext_len", len(plaintext)))
	return string(plaintext), nil
}
