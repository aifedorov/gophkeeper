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

func (c *binaryClient) Upload(ctx context.Context, fileInfo *binary.FileInfo, reader io.Reader) (id string, version int64, err error) {
	clientStream, err := c.client.Upload(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create upload stream: %w", err)
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
		return "", 0, fmt.Errorf("failed to send file metadata: %w", err)
	}

	buffer := make([]byte, bufferSize)
	var uploaded int64
	for {
		n, err := reader.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", 0, fmt.Errorf("failed to read file: %w", err)
		}

		err = clientStream.Send(
			&pb.UploadRequest{
				Data: &pb.UploadRequest_Chunk{
					Chunk: buffer[:n],
				},
			})
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, recvErr := clientStream.CloseAndRecv()
				if recvErr != nil {
					return "", 0, handleGRPCError(recvErr)
				}
			}
			return "", 0, fmt.Errorf("failed to send chunk: %w", err)
		}

		uploaded += int64(n)
		if size > 0 {
			fmt.Printf("\rUploaded: %d / %d bytes (%.1f%%)", uploaded, size, float64(uploaded)/float64(size)*100)
		}
	}
	if size > 0 {
		fmt.Println()
	}

	resp, err := clientStream.CloseAndRecv()
	if err != nil {
		return "", 0, handleGRPCError(err)
	}

	return resp.GetFileId(), resp.GetVersion(), nil
}

func (c *binaryClient) List(ctx context.Context) ([]binary.File, error) {
	response, err := c.client.List(ctx, &pb.ListRequest{})
	if err != nil {
		return nil, handleGRPCError(err)
	}

	files := make([]binary.File, len(response.GetFiles()))
	for i, file := range response.GetFiles() {
		dFile, err := binary.NewFile(
			file.GetId(),
			file.GetName(),
			file.GetSize(),
			file.GetNotes(),
			file.GetVersion(),
			file.GetUploadedAt().AsTime(),
		)
		if err != nil || dFile == nil {
			return nil, fmt.Errorf("failed to convert file metadata: %w", err)
		}
		files[i] = *dFile
	}
	return files, nil
}

func (c *binaryClient) Download(ctx context.Context, id string) (io.ReadCloser, *binary.FileMeta, error) {
	req := pb.DownloadRequest{
		FileId: &id,
	}
	stream, err := c.client.Download(ctx, &req)
	if err != nil {
		return nil, nil, handleGRPCError(err)
	}

	firstMsg, err := stream.Recv()
	if err != nil {
		return nil, nil, handleGRPCError(err)
	}
	meta := firstMsg.GetFile()
	fileMeta, err := binary.NewFileMeta(
		meta.GetName(),
		meta.GetSize(),
		meta.GetNotes(),
		meta.GetVersion(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file metadata: %w", err)
	}

	return newGRPCStreamReader(stream), fileMeta, nil
}

func (c *binaryClient) Update(ctx context.Context, fileInfo *binary.UpdateFileInfo, reader io.Reader) (version int64, err error) {
	clientStream, err := c.client.Update(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create update stream: %w", err)
	}

	id := fileInfo.ID()
	name := fileInfo.Name()
	size := fileInfo.Size()
	notes := fileInfo.Notes()
	ver := fileInfo.Version()
	req := pb.UpdateRequest{
		Data: &pb.UpdateRequest_File{
			File: &pb.UpdateRequest_Metadata{
				FileId:  &id,
				Name:    &name,
				Size:    &size,
				Notes:   &notes,
				Version: &ver,
			},
		},
	}
	err = clientStream.Send(&req)
	if err != nil {
		return 0, fmt.Errorf("failed to send file metadata: %w", err)
	}

	buffer := make([]byte, bufferSize)
	var uploaded int64
	for {
		n, err := reader.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to read file: %w", err)
		}

		err = clientStream.Send(
			&pb.UpdateRequest{
				Data: &pb.UpdateRequest_Chunk{
					Chunk: buffer[:n],
				},
			})
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, recvErr := clientStream.CloseAndRecv()
				if recvErr != nil {
					return 0, handleGRPCError(recvErr)
				}
			}
			return 0, fmt.Errorf("failed to send chunk: %w", err)
		}

		uploaded += int64(n)
		if size > 0 {
			fmt.Printf("\rUploaded: %d / %d bytes (%.1f%%)", uploaded, size, float64(uploaded)/float64(size)*100)
		}
	}
	if size > 0 {
		fmt.Println()
	}

	resp, err := clientStream.CloseAndRecv()
	if err != nil {
		return 0, handleGRPCError(err)
	}

	return resp.GetVersion(), nil
}

func (c *binaryClient) Delete(ctx context.Context, id string) error {
	req := pb.DeleteRequest{
		FileId: &id,
	}
	_, err := c.client.Delete(ctx, &req)
	if err != nil {
		return handleGRPCError(err)
	}
	return nil
}
