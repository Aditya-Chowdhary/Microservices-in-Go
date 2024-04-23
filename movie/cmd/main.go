package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"movie-micro/gen"
	"movie-micro/movie/internal/controller/movie"
	metadatagateway "movie-micro/movie/internal/gateway/metadata/grpc"
	ratinggateway "movie-micro/movie/internal/gateway/rating/grpc"
	grpchandler "movie-micro/movie/internal/handler/grpc"
	"movie-micro/pkg/discovery"
	"movie-micro/pkg/discovery/consul"
	"movie-micro/pkg/tracing"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
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
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	f, err := os.Open("./configs/base.yaml")
	if err != nil {
		logger.Fatal("Failed to open configuration", zap.Error(err))
	}
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		logger.Fatal("Failed to parse configuration", zap.Error(err))
	}
	port := cfg.API.Port

	logger.Info("Starting the movie service", zap.Int("port", port))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := tracing.NewJaegerProvider(cfg.Jaeger.URL,
		serviceName)
	if err != nil {
		logger.Fatal("Failed to initialize Jaeger provider", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Fatal("Failed to shut down Jaeger prodiver", zap.Error(err))
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// ! Code for using consul service registry
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		logger.Fatal("Failed to create consule registry on port 8500: ", zap.Error(err))
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		logger.Fatal("Failed to generate instanceID for metadata: ", zap.Error(err))
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				logger.Error("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	// ! Code for using memory service registry
	// registry := memory.NewRegistry()
	// ctx := context.Background()
	// metadatainstanceID := discovery.GenerateInstanceID(serviceName)
	// if err := registry.Register(ctx, metadatainstanceID, "movie", "localhost:8081"); err != nil {
	// 	panic(err)
	// }
	// ratinginstanceID := discovery.GenerateInstanceID(serviceName)
	// if err := registry.Register(ctx, ratinginstanceID, "movie", "localhost:8082"); err != nil {
	// 	panic(err)
	// }
	// movieinstanceID := discovery.GenerateInstanceID(serviceName)
	// if err := registry.Register(ctx, movieinstanceID, "movie", "localhost:8083"); err != nil {
	// 	panic(err)
	// }
	// defer registry.Deregister(ctx, movieinstanceID, "movie")

	// ! Unchanged
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// const limit = 100
	// const burst = 100
	// l := newLimiter(100, 100)
	// srv := grpc.NewServer(grpc.UnaryInterceptor(ratelimit.UnaryServerInterceptor(l)), grpc.StatsHandler(otelgrpc.NewServerHandler()))
	srv := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
