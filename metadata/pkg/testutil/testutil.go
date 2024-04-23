package testutil

import (
	"movie-micro/gen"
	"movie-micro/metadata/internal/controller/metadata"
	grpchandler "movie-micro/metadata/internal/handler/grpc"
	"movie-micro/metadata/internal/repository/memory"
)

// NewTestMetadataGRPCServer creates a new metadata gRPC server for tests
func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpchandler.New(ctrl)
}
