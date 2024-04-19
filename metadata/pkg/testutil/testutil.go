package testutil

import (
	"movie-micro/gen"
	"movie-micro/metadata/internal/controller/metadata"
	"movie-micro/metadata/internal/repository/memory"
	grpchandler "movie-micro/metadata/internal/handler/grpc"
)

// NewTestMetadataGRPCServer creates a new metadata gRPC server for tests
func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpchandler.New(ctrl)
}