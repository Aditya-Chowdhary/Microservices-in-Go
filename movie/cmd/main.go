package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"movie-micro/gen"
	"movie-micro/movie/internal/controller/movie"
	metadatagateway "movie-micro/movie/internal/gateway/metadata/http"
	ratinggateway "movie-micro/movie/internal/gateway/rating/http"
	grpchandler "movie-micro/movie/internal/handler/grpc"
	"movie-micro/pkg/discovery"
	"movie-micro/pkg/discovery/memory"

	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "movie"

type limiter struct {
	l *rate.Limiter
}

func newLimiter(limit, burst int) *limiter {
	return &limiter{rate.NewLimiter(rate.Limit(limit), burst)}
}

func (l *limiter) Limit() bool {
	return l.l.Allow()
}

func main() {
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port

	log.Printf("Starting the movie service on port %d", port)

	// ! Code for using consul service registry
	/*
		// registry, err := consul.NewRegistry("localhost:8500")
		// if err != nil {
		// 	panic(err)
		// }

		// ctx := context.Background()
		// instanceID := discovery.GenerateInstanceID(serviceName)
		// if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		// 	panic(err)
		// }

		// go func() {
		// 	for {
		// 		if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
		// 			log.Println("Failed to report healthy state: " + err.Error())
		// 		}
		// 		time.Sleep(1 * time.Second)
		// 	}
		// }()
		// defer registry.Deregister(ctx, instanceID, serviceName) */

	// ! Code for using memory service registry
	registry := memory.NewRegistry()
	ctx := context.Background()
	metadatainstanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, metadatainstanceID, "movie", "localhost:8081"); err != nil {
		panic(err)
	}
	ratinginstanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, ratinginstanceID, "movie", "localhost:8082"); err != nil {
		panic(err)
	}
	movieinstanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, movieinstanceID, "movie", "localhost:8083"); err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, movieinstanceID, "movie")

	// ! Unchanged
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	const limit = 100
	const burst = 100
	l := newLimiter(100, 100)
	srv := grpc.NewServer(grpc.UnaryInterceptor(ratelimit.UnaryServerInterceptor(l)))
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
