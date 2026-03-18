package client

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/credential/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type credentialClient struct {
	client pb.CredentialServiceClient
}

func NewCredentialClient(conn *grpc.ClientConn) credential.Client {
	return &credentialClient{
		client: pb.NewCredentialServiceClient(conn),
	}
}

func (c *credentialClient) Create(ctx context.Context, creds credential.Credential) (id string, version int64, err error) {
	request := pb.CreateRequest{
		Name:     &creds.Name,
		Login:    &creds.Login,
		Password: &creds.Password,
		Notes:    &creds.Notes,
	}

	res, err := c.client.Create(ctx, &request)
	if err != nil {
		return "", 0, handleGRPCError(err)
	}

	return res.GetId(), res.GetVersion(), nil
}

func (c *credentialClient) Update(ctx context.Context, id string, creds credential.Credential) (version int64, err error) {
	request := pb.UpdateRequest{
		Id:       &id,
		Name:     &creds.Name,
		Login:    &creds.Login,
		Password: &creds.Password,
		Notes:    &creds.Notes,
		Version:  &creds.Version,
	}
	response, err := c.client.Update(ctx, &request)
	if err != nil {
		return 0, handleGRPCError(err)
	}

	return response.GetVersion(), nil
}

func (c *credentialClient) Delete(ctx context.Context, id string) error {
	request := pb.DeleteRequest{
		Id: &id,
	}

	resp, err := c.client.Delete(ctx, &request)
	if err != nil {
		return handleGRPCError(err)
	}
	if !resp.GetSuccess() {
		return fmt.Errorf("client: delete operation failed")
	}
	return nil
}

func (c *credentialClient) List(ctx context.Context) ([]credential.Credential, error) {
	request := pb.ListRequest{}
	response, err := c.client.List(ctx, &request)
	if err != nil {
		return []credential.Credential{}, handleGRPCError(err)
	}
	credentials := make([]credential.Credential, len(response.Credentials))
	for i, cred := range response.Credentials {
		credentials[i] = credential.Credential{
			ID:       cred.GetId(),
			Name:     cred.GetName(),
			Login:    cred.GetLogin(),
			Password: cred.GetPassword(),
			Notes:    cred.GetNotes(),
			Version:  cred.GetVersion(),
		}
	}
	return credentials, nil
}

// handleGRPCError extracts gRPC status codes from errors and maps them to domain-specific errors.
// This allows the rest of the application to handle errors without knowing about gRPC implementation details.
func handleGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("client: failed to extract gRPC status: %w", err)
	}

	switch st.Code() {
	case codes.Aborted:
		return shared.ErrVersionConflict
	case codes.NotFound:
		return shared.ErrNotFound
	case codes.AlreadyExists:
		return shared.ErrAlreadyExists
	case codes.Unauthenticated:
		return shared.ErrUnauthenticated
	default:
		return fmt.Errorf("client: unexpected gRPC error: %w", err)
	}
}
