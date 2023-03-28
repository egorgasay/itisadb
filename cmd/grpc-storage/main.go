package main

import (
	"fmt"
	"github.com/egorgasay/grpc-storage/config"
	"github.com/egorgasay/grpc-storage/internal/handler"
	"github.com/egorgasay/grpc-storage/internal/storage"
	"github.com/egorgasay/grpc-storage/internal/usecase"
	pb "github.com/egorgasay/grpc-storage/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()

	store, err := storage.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	logic := usecase.New(store)
	h := handler.New(logic)
	grpcServer := grpc.NewServer()

	go func() {
		log.Println("Starting Server ...")
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:80"))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		pb.RegisterStorageServer(grpcServer, h)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logic.Save()
	log.Println("Shutdown Server ...")
}
