package client

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/binary/v1"
	"google.golang.org/grpc"
)

const bufferSize = 1024 * 1024

type binaryClient struct {
	client pb.BinaryServiceClient
}

func NewBinaryClient(conn *grpc.ClientConn) binary.Client {
	return &binaryClient{
		client: pb.NewBinaryServiceClient(conn),
	}
}

func (c *binaryClient) Upload(ctx context.Context, fileInfo *binary.FileInfo, reader io.Reader) error {
	clientStream, err := c.client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("failed to create upload stream: %w", err)
	}

	name := fileInfo.Name()
	size := fileInfo.Size()
	mimeType := fileInfo.MimeType()
	req := pb.UploadRequest{
		Data: &pb.UploadRequest_File{
			File: &pb.UploadRequest_FileMetadata{
				Name:     &name,
				Size:     &size,
				MimeType: &mimeType,
			},
		},
	}
	err = clientStream.Send(&req)
	if err != nil {
		return fmt.Errorf("failed to send file metadata: %w", err)
	}

	buffer := make([]byte, bufferSize)
	for {
		n, err := reader.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		err = clientStream.Send(
			&pb.UploadRequest{
				Data: &pb.UploadRequest_Chunk{
					Chunk: buffer[:n],
				},
			})
		if err != nil {
			return fmt.Errorf("failed to send chunk: %w", err)
		}
	}

	_, err = clientStream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to complete upload: %w", err)
	}

	return nil
}
