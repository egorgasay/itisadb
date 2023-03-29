package main

import (
	"fmt"
	"github.com/egorgasay/grpc-storage/internal/grpc-storage/config"
	"github.com/egorgasay/grpc-storage/internal/grpc-storage/handler"
	"github.com/egorgasay/grpc-storage/internal/grpc-storage/storage"
	"github.com/egorgasay/grpc-storage/internal/grpc-storage/usecase"
	"github.com/egorgasay/grpc-storage/pkg/api/balancer"
	pb "github.com/egorgasay/grpc-storage/pkg/api/storage"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		log.Fatalf("Failed to initialize: %v", err)
	}

	logic := usecase.New(store)
	h := handler.New(logic)

	conn, err := grpc.Dial(cfg.Balancer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	ram := usecase.RAMUsage()

	cr := &balancer.ConnectRequest{
		Address:   cfg.Host,
		Total:     ram.Total,
		Available: ram.Available,
	}

	cl := balancer.NewBalancerClient(conn)
	_, err = cl.Connect(context.Background(), cr)
	if err != nil {
		log.Fatalf("Unable to connect to the balancer: %v", err)
	}

	grpcServer := grpc.NewServer()

	go func() {
		log.Println("Starting Server ...")
		lis, err := net.Listen("tcp", fmt.Sprintf(cfg.Host))
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
