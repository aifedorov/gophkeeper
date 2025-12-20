package client

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/binary/v1"
)

func toDomain(metadata *pb.MetadataResponse) (*binary.File, error) {
	if metadata == nil {
		return nil, fmt.Errorf("failed to convert nil metadata")
	}
	return binary.NewFile(
		metadata.GetId(),
		metadata.GetName(),
		metadata.GetSize(),
		metadata.GetNotes(),
		metadata.GetUploadedAt().AsTime(),
	)
}
