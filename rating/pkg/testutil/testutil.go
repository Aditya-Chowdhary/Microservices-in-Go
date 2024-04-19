package testutil

import (
	"movie-micro/gen"
	"movie-micro/rating/internal/controller/rating"
	"movie-micro/rating/internal/repository/memory"
	grpchandler "movie-micro/rating/internal/handler/grpc"
)

// NewTestRatingGRPCServer creates a new rating gRPC server for tests
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}