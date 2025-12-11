package credential

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"github.com/google/uuid"
)

// TODO: Implement encryption/decryption logic for credentials
// Currently using direct byte conversion as placeholder
// Need to add proper encryption when storing and decryption when retrieving

func toDomainCredential(credential interfaces.RepositoryCredential) (Credential, error) {
	id, err := uuid.Parse(credential.ID)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to parse credential id: %w", err)
	}
	userID, err := uuid.Parse(credential.UserID)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to parse user id: %w", err)
	}
	return Credential{
		id:       id,
		userID:   userID,
		name:     credential.Name,
		login:    string(credential.Encryptedlogin),    // TODO: decrypt
		password: string(credential.Encryptedpassword), // TODO: decrypt
		metadata: string(credential.Encryptednotes),    // TODO: decrypt
	}, nil
}

func toRepositoryCredential(credential Credential) interfaces.RepositoryCredential {
	return interfaces.RepositoryCredential{
		ID:                credential.GetID().String(),
		UserID:            credential.GetUserID().String(),
		Name:              credential.GetName(),
		Encryptedlogin:    []byte(credential.GetLogin()),    // TODO: encrypt
		Encryptedpassword: []byte(credential.GetPassword()), // TODO: encrypt
		Encryptednotes:    []byte(credential.GetMetadata()), // TODO: encrypt
	}
}
