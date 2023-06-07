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
	"itisadb/pkg/logger"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	loggerInstance, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to inizialise logger: %v", err)
	}

	cfg := config.New()
	store, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("failed to inizialise logic layer: %v", err)
	}

	logic, err := usecase.New(store, logger.New(loggerInstance))
	if err != nil {
		log.Fatalf("failed to inizialise logic layer: %v", err)
	}

	h := grpchandler.New(logic)
	grpcServer := grpc.NewServer()

	log.Println("Starting Balancer ...")
	lis, err := net.Listen("tcp", cfg.GRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	balancer.RegisterBalancerServer(grpcServer, h)

	// gRPC by default
	go func() {
		log.Println("Starting GRPC", cfg.GRPC)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	handler := resthandler.New(logic)

	if cfg.REST != "" {
		go func() {
			log.Printf("Starting FastHTTP %s", cfg.REST)
			if err := fasthttp.ListenAndServe(cfg.REST, handler.ServeHTTP); err != nil {
				log.Fatalf("Error in ListenAndServe: %v", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	grpcServer.GracefulStop()
	log.Println("Shutdown Balancer ...")
}
