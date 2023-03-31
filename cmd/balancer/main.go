package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"grpc-storage/internal/memory-balancer/config"
	"grpc-storage/internal/memory-balancer/handler"
	"grpc-storage/internal/memory-balancer/storage"
	"grpc-storage/internal/memory-balancer/usecase"
	balancer "grpc-storage/pkg/api/balancer"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to inizialise logger: %v", err)
	}

	cfg := config.New()
	store, err := storage.New(cfg)
	logic := usecase.New(store, logger)
	h := handler.New(logic)
	grpcServer := grpc.NewServer()

	log.Println("Starting Balancer ...")
	lis, err := net.Listen("tcp", fmt.Sprintf(cfg.Host))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	balancer.RegisterBalancerServer(grpcServer, h)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	grpcServer.GracefulStop()
	log.Println("Shutdown Balancer ...")
}
