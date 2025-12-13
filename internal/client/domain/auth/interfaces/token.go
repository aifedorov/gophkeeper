package interfaces

import "context"

type SessionProvider interface {
	GetSession(ctx context.Context) (Session, error)
}
