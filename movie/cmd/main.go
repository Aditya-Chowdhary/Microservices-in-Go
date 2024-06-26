package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Aditya-Chowdhary/micro-movies/gen"
	"github.com/Aditya-Chowdhary/micro-movies/movie/internal/controller/movie"
	metadatagateway "github.com/Aditya-Chowdhary/micro-movies/movie/internal/gateway/metadata/grpc"
	ratinggateway "github.com/Aditya-Chowdhary/micro-movies/movie/internal/gateway/rating/grpc"
	grpchandler "github.com/Aditya-Chowdhary/micro-movies/movie/internal/handler/grpc"
	"github.com/Aditya-Chowdhary/micro-movies/pkg/discovery"
	"github.com/Aditya-Chowdhary/micro-movies/pkg/discovery/consul"
	"github.com/Aditya-Chowdhary/micro-movies/pkg/tracing"

	"github.com/uber-go/tally"
	"github.com/uber-go/tally/prometheus"
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

	// ! Code for using Jaeger Tracing
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

	// ! Code for prometheus reporting
	reporter := prometheus.NewReporter(prometheus.Options{})
	scope, closer := tally.NewRootScope(tally.ScopeOptions{
		Tags:                   map[string]string{"service": "movie"},
		CachedReporter:         reporter,
		OmitCardinalityMetrics: true,
	}, 10*time.Second)
	defer closer.Close()

	http.Handle("/metrics", reporter.HTTPHandler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Prometheus.MetricsPort), nil); err != nil {
			logger.Fatal("Failed to start the metrics handler", zap.Error(err))
		}
	}()

	counter := scope.Tagged(map[string]string{
		"service": "movie",
	}).Counter("service_started")
	counter.Inc(1)

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
