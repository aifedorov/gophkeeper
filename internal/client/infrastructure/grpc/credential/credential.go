package client

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/credential/v1"
	"google.golang.org/grpc"
)

type credentialClient struct {
	client pb.CredentialServiceClient
}

func NewCredentialClient(conn *grpc.ClientConn) credential.CredentialClient {
	return &credentialClient{
		client: pb.NewCredentialServiceClient(conn),
	}
}

func (c *credentialClient) Create(ctx context.Context, creds credential.Credential) error {
	request := pb.CreateRequest{
		Name:     &creds.Name,
		Login:    &creds.Login,
		Password: &creds.Password,
		Notes:    &creds.Notes,
	}
	_, err := c.client.Create(ctx, &request)
	if err != nil {
		return fmt.Errorf("client: failed to create credential: %w", err)
	}
	return nil
}

func (c *credentialClient) Get(ctx context.Context, id string) (credential.Credential, error) {
	request := pb.GetRequest{
		Id: &id,
	}
	response, err := c.client.Get(ctx, &request)
	if err != nil {
		return credential.Credential{}, fmt.Errorf("client: failed to get credential: %w", err)
	}
	return credential.Credential{
		ID:       response.GetId(),
		Name:     response.GetName(),
		Login:    response.GetLogin(),
		Password: response.GetPassword(),
		Notes:    response.GetNotes(),
	}, nil
}

func (c *credentialClient) Update(ctx context.Context, id string, creds credential.Credential) error {
	request := pb.UpdateRequest{
		Id:       &id,
		Name:     &creds.Name,
		Login:    &creds.Login,
		Password: &creds.Password,
		Notes:    &creds.Notes,
	}
	_, err := c.client.Update(ctx, &request)
	if err != nil {
		return fmt.Errorf("client: failed to update credential: %w", err)
	}
	return nil
}

func (c *credentialClient) Delete(ctx context.Context, id string) error {
	request := pb.DeleteRequest{
		Id: &id,
	}
	_, err := c.client.Delete(ctx, &request)
	if err != nil {
		return fmt.Errorf("client: failed to delete credential: %w", err)
	}
	return nil
}

func (c *credentialClient) List(ctx context.Context) ([]credential.Credential, error) {
	request := pb.ListRequest{}
	response, err := c.client.List(ctx, &request)
	if err != nil {
		return []credential.Credential{}, fmt.Errorf("client: failed to list credentials: %w", err)
	}
	credentials := make([]credential.Credential, len(response.Credentials))
	for i, cred := range response.Credentials {
		credentials[i] = credential.Credential{
			ID:       cred.GetId(),
			Name:     cred.GetName(),
			Login:    cred.GetLogin(),
			Password: cred.GetPassword(),
			Notes:    cred.GetNotes(),
		}
	}
	return credentials, nil
}
