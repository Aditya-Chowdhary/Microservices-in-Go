package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"movie-micro/gen"
	"movie-micro/rating/internal/controller/rating"
	grpchandler "movie-micro/rating/internal/handler/grpc"
	"movie-micro/rating/internal/repository/mysql"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "rating"

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

	log.Printf("Starting the rating metadata service on port %d", port)

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

	_, cancel := context.WithCancel(context.Background())

	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}
	ctrl := rating.New(repo, nil)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := <-sigChan
		cancel()
		log.Printf("Received signal %v, attempting graceful shutdown\n", s)
		srv.GracefulStop()
		log.Println("Gracefully stopped the gRPC server")
	}()
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
	wg.Wait()
}
