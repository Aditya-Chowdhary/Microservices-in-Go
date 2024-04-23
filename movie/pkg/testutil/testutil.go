package testutil

import (
	"github.com/Aditya-Chowdhary/micro-movies/gen"
	"github.com/Aditya-Chowdhary/micro-movies/movie/internal/controller/movie"
	metadatagateway "github.com/Aditya-Chowdhary/micro-movies/movie/internal/gateway/metadata/grpc"
	ratinggateway "github.com/Aditya-Chowdhary/micro-movies/movie/internal/gateway/rating/grpc"
	grpchandler "github.com/Aditya-Chowdhary/micro-movies/movie/internal/handler/grpc"
	"github.com/Aditya-Chowdhary/micro-movies/pkg/discovery"
)

// NewTestMovieGRPCServer creates a new movie gRPC server for tests
func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	return grpchandler.New(ctrl)
}
