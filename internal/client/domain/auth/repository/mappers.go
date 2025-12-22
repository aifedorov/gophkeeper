package repository

import (
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage"
)

func toDomainSession(session storage.Session) interfaces.Session {
	return interfaces.NewSession(
		session.GetAccessToken(),
		session.GetEncryptionKey(),
		session.GetUserID(),
		session.GetLogin(),
	)
}

func toStoreSession(session interfaces.Session) storage.Session {
	return storage.NewSession(
		session.GetAccessToken(),
		session.GetEncryptionKey(),
		session.GetUserID(),
		session.GetLogin(),
	)
}
