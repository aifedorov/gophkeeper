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
	ctx      context.Context
	mStorage *storage.Storage
}

func NewRepository(ctx context.Context, mStorage *storage.Storage) Repository {
	return &repository{
		ctx:      ctx,
		mStorage: mStorage,
	}
}

func (r *repository) Save(session interfaces.Session) error {
	return r.mStorage.Save(toMemorySession(session))
}

func (r *repository) Load() (interfaces.Session, error) {
	session, err := r.mStorage.Load()
	if err != nil {
		return interfaces.Session{}, auth.ErrSessionNotFound
	}
	return toDomainSession(session), nil
}

func (r *repository) Delete() error {
	return r.mStorage.Delete()
}

func toDomainUser(user storage.User) interfaces.User {
	return interfaces.User{
		ID:    user.ID,
		Login: user.Login,
	}
}

func toDomainSession(session storage.Session) interfaces.Session {
	return interfaces.Session{
		User:        toDomainUser(session.User),
		AccessToken: session.AccessToken,
	}
}

func toMemorySession(session interfaces.Session) storage.Session {
	return storage.Session{
		User:        toMemoryUser(session.User),
		AccessToken: session.AccessToken,
	}
}

func toMemoryUser(user interfaces.User) storage.User {
	return storage.User{
		ID:    user.ID,
		Login: user.Login,
	}
}
