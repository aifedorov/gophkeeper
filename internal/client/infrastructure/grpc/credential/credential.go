package grpc

import (
	"context"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/credential/v1"
	"google.golang.org/grpc"
)

type CredentialClient interface {
	Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error)
	Get(context.Context, *pb.GetRequest) (*pb.GetResponse, error)
	Update(context.Context, *pb.UpdateRequest) (*pb.UpdateResponse, error)
	Delete(context.Context, *pb.DeleteRequest) (*pb.DeleteResponse, error)
	List(context.Context, *pb.ListRequest) (*pb.ListResponse, error)
}

type credentialClient struct {
	client pb.CredentialServiceClient
}

func NewCredentialClient(conn *grpc.ClientConn) CredentialClient {
	return &credentialClient{
		client: pb.NewCredentialServiceClient(conn),
	}
}

func (c *credentialClient) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	return c.client.Create(ctx, request)
}

func (c *credentialClient) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *credentialClient) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *credentialClient) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *credentialClient) List(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {
	//TODO implement me
	panic("implement me")
}
