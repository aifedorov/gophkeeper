package interfaces

import "context"

type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}
