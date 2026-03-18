package server

import (
	"errors"
	"fmt"
	"io"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/binary/v1"
	"google.golang.org/grpc"
)

type uploadStreamReader struct {
	stream grpc.ClientStreamingServer[pb.UploadRequest, pb.UploadResponse]
	buffer []byte
}

func newUploadStreamReader(stream grpc.ClientStreamingServer[pb.UploadRequest, pb.UploadResponse]) *uploadStreamReader {
	return &uploadStreamReader{
		stream: stream,
		buffer: nil,
	}
}

func (r *uploadStreamReader) Read(p []byte) (n int, err error) {
	if len(r.buffer) > 0 {
		n = copy(p, r.buffer)
		r.buffer = r.buffer[n:]
		return n, nil
	}

	msg, err := r.stream.Recv()
	if errors.Is(err, io.EOF) {
		return 0, io.EOF
	}
	if err != nil {
		return 0, fmt.Errorf("failed to receive message: %w", err)
	}

	chunk := msg.GetChunk()
	n = copy(p, chunk)

	if n < len(chunk) {
		r.buffer = chunk[n:]
	}

	return n, nil
}
