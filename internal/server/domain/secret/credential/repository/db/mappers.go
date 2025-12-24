package repository

import "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"

// toInterfacesCredential converts a database Credential model to the repository interface representation.
// This mapper transforms sqlc-generated types to domain-layer types.
func toInterfacesCredential(c Credential) interfaces.RepositoryCredential {
	return interfaces.RepositoryCredential{
		ID:                c.ID.String(),
		UserID:            c.UserID.String(),
		Name:              c.Name,
		EncryptedLogin:    c.Encryptedlogin,
		EncryptedPassword: c.Encryptedpassword,
		EncryptedNotes:    c.Encryptednotes,
		Version:           c.Version,
	}
}
