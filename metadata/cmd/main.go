package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"movie-micro/gen"
	"movie-micro/metadata/internal/controller/metadata"
	grpchandler "movie-micro/metadata/internal/handler/grpc"
	"movie-micro/metadata/internal/repository/memory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "metadata"

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

	log.Printf("Starting the movie metadata service on port %d", port)

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
	// defer registry.Deregister(ctx, instanceID, serviceName)
	*/

	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
