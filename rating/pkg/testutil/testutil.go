package testutil

import (
	"github.com/Aditya-Chowdhary/micro-movies/gen"
	"github.com/Aditya-Chowdhary/micro-movies/rating/internal/controller/rating"
	grpchandler "github.com/Aditya-Chowdhary/micro-movies/rating/internal/handler/grpc"
	"github.com/Aditya-Chowdhary/micro-movies/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server for tests
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}
