package main

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"itisadb/internal/memory-balancer/config"
	grpchandler "itisadb/internal/memory-balancer/handler/grpc"
	resthandler "itisadb/internal/memory-balancer/handler/rest"
	"itisadb/internal/memory-balancer/storage"
	"itisadb/internal/memory-balancer/usecase"
	balancer "itisadb/pkg/api/balancer"
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
	if err != nil {
		log.Fatalf("failed to inizialise logic layer: %v", err)
	}

	logic, err := usecase.New(store, logger)
	if err != nil {
		log.Fatalf("failed to inizialise logic layer: %v", err)
	}

	h := grpchandler.New(logic)
	grpcServer := grpc.NewServer()

	log.Println("Starting Balancer ...")
	lis, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	balancer.RegisterBalancerServer(grpcServer, h)

	go func() {
		log.Println("Starting GRPC", cfg.Host)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	handler := resthandler.New(logic)

	go func() {
		log.Println("Starting FastHTTP 127.0.0.1:890")
		if err := fasthttp.ListenAndServe("127.0.0.1:890", handler.ServeHTTP); err != nil {
			log.Fatalf("Error in ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	grpcServer.GracefulStop()
	log.Println("Shutdown Balancer ...")
}
