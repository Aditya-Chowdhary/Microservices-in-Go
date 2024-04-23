package testutil

import (
	"movie-micro/gen"
	"movie-micro/rating/internal/controller/rating"
	grpchandler "movie-micro/rating/internal/handler/grpc"
	"movie-micro/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server for tests
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}
