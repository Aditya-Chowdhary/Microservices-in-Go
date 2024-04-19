package testutil

import (
	"movie-micro/gen"
	"movie-micro/movie/internal/controller/movie"
	metadatagateway "movie-micro/movie/internal/gateway/metadata/grpc"
	ratinggateway "movie-micro/movie/internal/gateway/rating/grpc"
	grpchandler "movie-micro/movie/internal/handler/grpc"
	"movie-micro/pkg/discovery"
)

// NewTestMovieGRPCServer creates a new movie gRPC server for tests
func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	return grpchandler.New(ctrl)
}
