package testutil

import (
	"github.com/Aditya-Chowdhary/micro-movies/gen"
	"github.com/Aditya-Chowdhary/micro-movies/metadata/internal/controller/metadata"
	grpchandler "github.com/Aditya-Chowdhary/micro-movies/metadata/internal/handler/grpc"
	"github.com/Aditya-Chowdhary/micro-movies/metadata/internal/repository/memory"
)

// NewTestMetadataGRPCServer creates a new metadata gRPC server for tests
func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpchandler.New(ctrl)
}
