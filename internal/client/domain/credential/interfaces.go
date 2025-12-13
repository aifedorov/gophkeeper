package credential

import (
	"context"
)

type CredentialClient interface {
	Create(ctx context.Context, creds Credential) error
	Get(ctx context.Context, id string) (Credential, error)
	Update(ctx context.Context, id string, creds Credential) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]Credential, error)
}
