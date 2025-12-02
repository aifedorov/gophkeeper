package auth

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/memory"
	"go.uber.org/zap"
)

type Repository interface {
	Save(session Session) error
	Load() (Session, error)
	Delete() error
}

type repository struct {
	ctx      context.Context
	logger   *zap.Logger
	mStorage *memory.Storage
}

func NewRepository(ctx context.Context, logger *zap.Logger, mStorage *memory.Storage) Repository {
	return &repository{
		ctx:      ctx,
		logger:   logger,
		mStorage: mStorage,
	}
}

func (r *repository) Save(session Session) error {
	return r.mStorage.Save(toMemorySession(session))
}

func (r *repository) Load() (Session, error) {
	session, err := r.mStorage.Load()
	if err != nil {
		return Session{}, ErrSessionNotFound
	}
	return toDomainSession(session), nil
}

func (r *repository) Delete() error {
	return r.mStorage.Delete()
}

func toDomainUser(user memory.User) User {
	return User{
		ID:    user.ID,
		Login: user.Login,
	}
}

func toDomainSession(session memory.Session) Session {
	return Session{
		User:        toDomainUser(session.User),
		AccessToken: session.AccessToken,
	}
}

func toMemorySession(session Session) memory.Session {
	return memory.Session{
		User:        toMemoryUser(session.User),
		AccessToken: session.AccessToken,
	}
}

func toMemoryUser(user User) memory.User {
	return memory.User{
		ID:    user.ID,
		Login: user.Login,
	}
}
