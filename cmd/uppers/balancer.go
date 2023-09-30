package uppers

import (
	"context"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"itisadb/internal/config"
	"itisadb/internal/core"
	grpchandler "itisadb/internal/handler/grpc"
	resthandler "itisadb/internal/handler/rest"
	"itisadb/internal/storage"
	"itisadb/pkg"
	"itisadb/pkg/api/balancer"
	"itisadb/pkg/logger"
	"log"
	"net"
)

func UpBalancer(ctx context.Context) {
	loggerInstance, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to inizialise logger: %v", err)
	}
	loggerInstance = loggerInstance.Named("BALANCER")

	cfg := config.New()
	store, err := storage.New()
	if err != nil {
		log.Fatalf("failed to inizialise logic layer: %v", err)
	}

	logic, err := core.New(ctx, store, logger.New(loggerInstance))
	if err != nil {
		log.Fatalf("failed to inizialise logic layer: %v", err)
	}

	loggerInstance.Info("Starting Balancer ...")

	// gRPC by default
	go runBalancerGRPC(ctx, loggerInstance, logic, cfg)

	if cfg.REST != "" {
		go runBalancerREST(ctx, loggerInstance, logic, cfg)
	}
}

func runBalancerGRPC(ctx context.Context, l *zap.Logger, logic *core.Core, cfg *config.Config) {
	h := grpchandler.New(logic)
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", cfg.GRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	balancer.RegisterBalancerServer(grpcServer, h)

	err = pkg.WithContext(ctx, func() error {
		l.Info("Starting GRPC", zap.String("address", cfg.GRPC))
		err = grpcServer.Serve(lis)
		if err != nil {
			return fmt.Errorf("grpcServer Serve: %v", err)
		}

		return nil
	}, make(chan struct{}, 1), func() {
		grpcServer.GracefulStop()
	})
	l.Info("Shutdown GRPC Balancer ...")

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Error in GRPC: %v", err)
	}
}

func runBalancerREST(ctx context.Context, l *zap.Logger, logic *core.Core, cfg *config.Config) {
	handler := resthandler.New(logic)
	lis, err := net.Listen("tcp", cfg.REST)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = pkg.WithContext(ctx, func() error {
		l.Info("Starting FastHTTP %s", zap.String("address", cfg.REST))
		if err := fasthttp.Serve(lis, handler.ServeHTTP); err != nil {
			return fmt.Errorf("error in REST Serve: %w", err)
		}

		return nil
	}, make(chan struct{}, 1), func() {
		if err := lis.Close(); err != nil {
			l.Warn("Failed to close listener", zap.Error(err))
		}
	})
	l.Info("Shutdown REST Balancer ...")

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Error in REST: %v", err)
	}
}
