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
	notes := fileInfo.Notes()
	req := pb.UploadRequest{
		Data: &pb.UploadRequest_File{
			File: &pb.UploadRequest_Metadata{
				Name:  &name,
				Size:  &size,
				Notes: &notes,
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

func (c *binaryClient) List(ctx context.Context) ([]binary.File, error) {
	response, err := c.client.List(ctx, &pb.ListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]binary.File, len(response.GetFiles()))
	for i, file := range response.GetFiles() {
		domainFile, err := toDomain(file)
		if err != nil || domainFile == nil {
			return nil, fmt.Errorf("failed to convert file metadata: %w", err)
		}
		files[i] = *domainFile
	}
	return files, nil
}

func (c *binaryClient) Download(ctx context.Context, id string) (io.ReadCloser, *binary.FileMeta, error) {
	req := pb.DownloadRequest{
		FileId: &id,
	}
	stream, err := c.client.Download(ctx, &req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file: %w", err)
	}

	firstMsg, err := stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to receive file metadata: %w", err)
	}
	meta := firstMsg.GetFile()
	fileMeta, err := binary.NewFileMeta(meta.GetName(), meta.GetSize(), meta.GetNotes())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file metadata: %w", err)
	}
	return newGRPCStreamReader(stream), fileMeta, nil
}

func (c *binaryClient) Delete(ctx context.Context, id string) error {
	req := pb.DeleteRequest{
		FileId: &id,
	}
	_, err := c.client.Delete(ctx, &req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
