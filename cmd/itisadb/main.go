package main

import (
	"context"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/domains"
	"itisadb/internal/service/core"
	"itisadb/internal/service/generator"
	"itisadb/internal/service/servers"
	"itisadb/internal/service/session"
	transactionlogger "itisadb/internal/service/transaction-logger"
	"itisadb/internal/storage"
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

	cfg, err := config.New()
	if err != nil {
		lg.Fatal("failed to inizialise config: %v", zap.Error(err))
	}

	store, err := storage.New()
	if err != nil {
		lg.Fatal("failed to inizialise storage: %v", zap.String("error", err.Error()))
	}

	ctx, cancel := context.WithCancel(context.Background())
	s, err := servers.New()
	if err != nil {
		lg.Fatal("failed to inizialise servers: %v", zap.Error(err))
	}

	var tl domains.TransactionLogger
	if cfg.TransactionLogger.On {
		tl, err = transactionlogger.New()
		if err != nil {
			lg.Fatal("failed to inizialise transaction logger: %v", zap.Error(err))
		}

		lg.Info("Transaction logger enabled")

		lg.Info("Starting recovery from transaction logger")
		if err = tl.Restore(store); err != nil {
			lg.Fatal("failed to restore transaction logger: %v", zap.Error(err))
		}
		lg.Info("Transaction logger recovery completed")

		tl.Run()
		lg.Info("Transaction logger started")
	} else {
		lg.Info("Transaction logger disabled")
	}

	gen := generator.New(lg)
	ses := session.New(store, gen, lg)

	logic, err := core.New(ctx, cfg, lg, store, tl, s, ses)
	if err != nil {
		lg.Fatal("failed to inizialise logic layer: %v", zap.String("error", err.Error()))
	}

	go runGRPC(ctx, lg, logic, cfg.Network, ses)
	go runWebCLI(ctx, cfg.WebApp, lg, cfg.Network.GRPC)

	if cfg.Network.REST != "" {
		go runREST(ctx, lg, logic, cfg.Network)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cancel()
	time.Sleep(1 * time.Second)
}
