package repository

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage"
)

//go:generate mockgen -source=repository.go -destination=mock_repository_test.go -package=repository

type Repository interface {
	Save(session interfaces.Session) error
	Load() (interfaces.Session, error)
	Delete() error
}

type repository struct {
	ctx   context.Context
	store *storage.Storage
}

func NewRepository(ctx context.Context, store *storage.Storage) Repository {
	return &repository{
		ctx:   ctx,
		store: store,
	}
}

func (r *repository) Save(session interfaces.Session) error {
	return r.store.Save(toStoreSession(session))
}

func (r *repository) Load() (interfaces.Session, error) {
	session, err := r.store.Load()
	if err != nil {
		return interfaces.Session{}, auth.ErrSessionNotFound
	}
	return toDomainSession(session), nil
}

func (r *repository) Delete() error {
	return r.store.Delete()
}

func toDomainSession(session storage.Session) interfaces.Session {
	return interfaces.NewSession(
		session.GetAccessToken(),
		session.GetEncryptionKey(),
		session.GetUserID(),
	)
}

func toStoreSession(session interfaces.Session) storage.Session {
	return storage.NewSession(
		session.GetAccessToken(),
		session.GetEncryptionKey(),
		session.GetUserID(),
	)
}
