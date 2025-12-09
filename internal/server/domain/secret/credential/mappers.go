package credential

import (
	repository "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/repository/db"
)

func toDomainCredential(credential repository.Credential) Credential {
	return Credential{
		id:        credential.ID,
		name:      credential.Name,
		login:     credential.Login,
		password:  credential.Password,
		metadata:  credential.Metadata,
		updatedAt: credential.UpdatedAt.Time,
	}
}

func toRepositoryCredential(credential Credential) repository.Credential {
	return repository.Credential{
		ID:       credential.id,
		Name:     credential.name,
		Login:    credential.login,
		Password: credential.password,
		Metadata: credential.metadata,
	}
}
