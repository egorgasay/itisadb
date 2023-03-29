package main

import (
	"fmt"
	"github.com/egorgasay/grpc-storage/internal/memory-balancer/config"
	"github.com/egorgasay/grpc-storage/internal/memory-balancer/handler"
	"github.com/egorgasay/grpc-storage/internal/memory-balancer/usecase"
	balancer "github.com/egorgasay/grpc-storage/pkg/api/balancer"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()
	logic := usecase.New()
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
	log.Println("Shutdown Balancer ...")
}
