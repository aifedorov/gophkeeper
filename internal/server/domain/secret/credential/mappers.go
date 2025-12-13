package credential

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"github.com/google/uuid"
)

func toDomainCredential(crypto interfaces.CryptoService, key []byte, credential interfaces.RepositoryCredential) (Credential, error) {
	id, err := uuid.Parse(credential.ID)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to parse credential id: %w", err)
	}
	userID, err := uuid.Parse(credential.UserID)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to parse user id: %w", err)
	}
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
		id:       id,
		userID:   userID,
		name:     credential.Name,
		login:    login,
		password: password,
		notes:    notes,
	}, nil
}

func toRepositoryCredential(crypto interfaces.CryptoService, key []byte, credential Credential) (interfaces.RepositoryCredential, error) {
	encryptLogin, err := crypto.Encrypt(credential.login, key)
	if err != nil {
		return interfaces.RepositoryCredential{}, fmt.Errorf("failed to encrypt login: %w", err)
	}

	ecryptPassword, err := crypto.Encrypt(credential.password, key)
	if err != nil {
		return interfaces.RepositoryCredential{}, fmt.Errorf("failed to encrypt password: %w", err)
	}

	ecryptNotes, err := crypto.Encrypt(credential.notes, key)
	if err != nil {
		return interfaces.RepositoryCredential{}, fmt.Errorf("failed to encrypt notes: %w", err)
	}

	return interfaces.RepositoryCredential{
		ID:                credential.GetID().String(),
		UserID:            credential.GetUserID().String(),
		Name:              credential.GetName(),
		EncryptedLogin:    encryptLogin,
		EncryptedPassword: ecryptPassword,
		EncryptedNotes:    ecryptNotes,
	}, nil
}
