package main

import (
	"fmt"
	"github.com/go-chi/httplog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-storage/internal/grpc-storage/config"
	"grpc-storage/internal/grpc-storage/handler"
	"grpc-storage/internal/grpc-storage/storage"
	"grpc-storage/internal/grpc-storage/usecase"
	"grpc-storage/pkg/api/balancer"
	pb "grpc-storage/pkg/api/storage"
	"grpc-storage/pkg/logger"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()

	lg := httplog.NewLogger("grpc-storage", httplog.Options{
		Concise: true,
	})
	store, err := storage.New(cfg.DBConfig, logger.New(lg))
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	logic := usecase.New(store, logger.New(lg))
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
	resp, err := cl.Connect(context.Background(), cr)
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
	_, err = cl.Disconnect(context.Background(), &balancer.DisconnectRequest{ServerNumber: resp.GetServerNumber()})
	if err != nil {
		log.Println(err)
	}
	grpcServer.GracefulStop()
	logic.Save()
	log.Println("Shutdown Server ...")
}
