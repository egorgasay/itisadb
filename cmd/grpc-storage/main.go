package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	globconfig "itisadb/config"
	"itisadb/internal/grpc-storage/config"
	"itisadb/internal/grpc-storage/handler"
	"itisadb/internal/grpc-storage/servernumber"
	"itisadb/internal/grpc-storage/storage"
	"itisadb/internal/grpc-storage/usecase"
	"itisadb/pkg/api/balancer"
	pb "itisadb/pkg/api/storage"
	"itisadb/pkg/logger"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gc := globconfig.New()
	err := gc.FromTOML(globconfig.DefaultConfigPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	cfg := config.New(gc.Storage)

	loggerInstance, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to inizialise logger: %v", err)
	}

	store, err := storage.New(cfg, logger.New(loggerInstance))
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	logic, err := usecase.New(store, logger.New(loggerInstance), cfg.IsTLogger)
	if err != nil {
		log.Fatalf("Failed to initialize usecase layer: %v", err)
	}

	h := handler.New(logic)

	conn, err := grpc.Dial(cfg.Balancer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	ram := usecase.RAMUsage()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current directory: %v", err)
	}

	cr := &balancer.BalancerConnectRequest{
		Address:   cfg.Host,
		Total:     ram.Total,
		Available: ram.Available,
		Server:    servernumber.Get(dir),
	}

	cl := balancer.NewBalancerClient(conn)
	resp, err := cl.Connect(context.Background(), cr)
	if err != nil {
		log.Fatalf("Unable to connect to the balancer: %v", err)
	}

	if cr.Server == 0 {
		err = servernumber.Set(resp.ServerNumber)
		if err != nil {
			log.Fatal(err.Error())
		}
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
	sc := bufio.NewScanner(os.Stdin)

	fmt.Println("PRESS ENTER FOR RECONNECT")
	for sc.Scan() {
		fmt.Print("RECONNECTING ...\n")
		cr.Server = resp.ServerNumber
		_, err = cl.Connect(context.Background(), cr)
		if err != nil {
			log.Println("Unable to connect to the balancer: %w", err)
		}
		fmt.Print("PRESS ENTER FOR RECONNECT")
	}

	<-quit
	_, err = cl.Disconnect(context.Background(), &balancer.BalancerDisconnectRequest{ServerNumber: resp.GetServerNumber()})
	if err != nil {
		log.Println(err)
	}
	grpcServer.GracefulStop()
	log.Println("Shutdown Server ...")
}
