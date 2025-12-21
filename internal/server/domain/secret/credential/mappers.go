// Package credential provides mappers for converting between domain and repository representations.
package credential

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
)

// toDomainCredential converts a repository credential representation to a domain Credential entity.
// It decrypts all encrypted fields (login, password, notes) using the provided encryption key.
// Returns an error if decryption fails for any field.
func toDomainCredential(crypto interfaces.CryptoService, key []byte, credential interfaces.RepositoryCredential) (Credential, error) {
	login, err := crypto.Decrypt(credential.EncryptedLogin, key)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to decrypt login: %w", err)
	}
	password, err := crypto.Decrypt(credential.EncryptedPassword, key)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to decrypt password: %w", err)
	}
	notes, err := crypto.Decrypt(credential.EncryptedNotes, key)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to decrypt notes: %w", err)
	}

	return Credential{
		id:       credential.ID,
		userID:   credential.UserID,
		name:     credential.Name,
		login:    login,
		password: password,
		notes:    notes,
	}, nil
}

// toRepositoryCredential converts a domain Credential entity to a repository credential representation.
// It encrypts all sensitive fields (login, password, notes) using the provided encryption key.
// Returns an error if encryption fails for any field.
func toRepositoryCredential(crypto interfaces.CryptoService, key []byte, credential Credential) (interfaces.RepositoryCredential, error) {
	encryptLogin, err := crypto.Encrypt(credential.login, key)
	if err != nil {
		return interfaces.RepositoryCredential{}, fmt.Errorf("failed to encrypt login: %w", err)
	}

	encryptPassword, err := crypto.Encrypt(credential.password, key)
	if err != nil {
		return interfaces.RepositoryCredential{}, fmt.Errorf("failed to encrypt password: %w", err)
	}

	encryptNotes, err := crypto.Encrypt(credential.notes, key)
	if err != nil {
		return interfaces.RepositoryCredential{}, fmt.Errorf("failed to encrypt notes: %w", err)
	}

	return interfaces.RepositoryCredential{
		ID:                credential.GetID(),
		UserID:            credential.GetUserID(),
		Name:              credential.GetName(),
		EncryptedLogin:    encryptLogin,
		EncryptedPassword: encryptPassword,
		EncryptedNotes:    encryptNotes,
	}, nil
}
