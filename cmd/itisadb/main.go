package main

import (
	"context"
	"go.uber.org/zap"
	"itisadb/internal/config"
	"itisadb/internal/core"
	"itisadb/internal/storage"
	"itisadb/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	lg, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to inizialise logger: %v", err)
	}

	cfg := config.New()
	store, err := storage.New()
	if err != nil {
		lg.Fatal("failed to inizialise storage: %v", zap.String("error", err.Error()))
	}

	ctx, cancel := context.WithCancel(context.Background())

	logic, err := core.New(ctx, store, logger.New(lg))
	if err != nil {
		lg.Fatal("failed to inizialise logic layer: %v", zap.String("error", err.Error()))
	}

	go runGRPC(ctx, lg, logic, cfg)
	go runWebCLI(ctx, lg)

	if cfg.REST != "" {
		go runREST(ctx, lg, logic, cfg)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cancel()
	time.Sleep(1 * time.Second)
}
