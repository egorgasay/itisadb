package uppers

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	globconfig "itisadb/config"
	"itisadb/internal/grpc-storage/config"
	"itisadb/internal/grpc-storage/handler"
	"itisadb/internal/grpc-storage/servernumber"
	"itisadb/internal/grpc-storage/storage"
	"itisadb/internal/grpc-storage/usecase"
	"itisadb/pkg"
	"itisadb/pkg/api/balancer"
	pb "itisadb/pkg/api/storage"
	"itisadb/pkg/logger"
	"log"
	"net"
	"os"
)

func UpGRPCStorage(ctx context.Context) {
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

	loggerInstance = loggerInstance.Named("STORAGE")

	store, err := storage.New()
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	logic, err := usecase.New(store, logger.New(loggerInstance), cfg.IsTLogger)
	if err != nil {
		log.Fatalf("Failed to initialize core layer: %v", err)
	}

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

	snum, err := servernumber.Get(dir)
	if err != nil {
		log.Printf("Unable to get server number: %v\n", err)
	}

	cr := &balancer.BalancerConnectRequest{
		Address:   cfg.Host,
		Total:     ram.Total,
		Available: ram.Available,
		Server:    snum,
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

	go func() {
		runGRPCStorage(ctx, loggerInstance, logic, cfg)

		_, err = cl.Disconnect(context.Background(), &balancer.BalancerDisconnectRequest{ServerNumber: resp.GetServerNumber()})
		if err != nil && !errors.Is(err, grpc.ErrClientConnClosing) {
			log.Fatalf("Error in GRPC Disconnect: %v", err)
		}
	}()
}

func runGRPCStorage(ctx context.Context, l *zap.Logger, logic *usecase.UseCase, cfg *config.Config) {
	grpcServer := grpc.NewServer()

	h := handler.New(logic)

	l.Info("Starting Server ...")
	lis, err := net.Listen("tcp", fmt.Sprintf(cfg.Host))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pb.RegisterStorageServer(grpcServer, h)

	err = pkg.WithContext(ctx, func() error {
		l.Info("Starting GRPC %s", zap.String("address", cfg.Host))
		err = grpcServer.Serve(lis)
		if err != nil {
			return fmt.Errorf("grpcServer Serve: %w", err)
		}

		return nil
	}, make(chan struct{}, 1), func() {
		grpcServer.GracefulStop()
	})

	l.Info("Shutdown Server ...")

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Error in GRPC: %v", err)
	}
}
