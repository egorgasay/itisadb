package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/domains"
	"itisadb/internal/service/balancer"
	"itisadb/internal/service/generator"
	"itisadb/internal/service/logic"
	"itisadb/internal/service/security"
	"itisadb/internal/service/servers"
	"itisadb/internal/service/session"
	"itisadb/internal/service/syncer"
	transactionlogger "itisadb/internal/service/transaction-logger"
	"itisadb/internal/storage"
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
		lg.Fatal("failed to inizialise storage", zap.String("error", err.Error()))
	}

	var tl domains.TransactionLogger

	if cfg.TransactionLogger.On {
		tl, err = transactionlogger.New(cfg.TransactionLogger)
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

	ctx, cancel := context.WithCancel(context.Background())

	appCFG := *cfg

	gen := generator.New(lg)
	ses := session.New(appCFG, store, gen, lg)
	sec := security.NewSecurityService(cfg.Security)

	uc := logic.NewLogic(store, *cfg, tl, lg, sec)

	var local = gost.None[domains.Server]()
	if !cfg.Balancer.On || (cfg.Balancer.On && !cfg.Balancer.BalancerOnly) {
		ls := servers.NewLocalServer(uc)
		local = local.Some(ls)
	}

	s, err := servers.New(local, lg)
	if err != nil {
		lg.Fatal("failed to inizialise balancer: %v", zap.Error(err))
	}

	// TODO: make it configurable
	syncer := syncer.NewSyncer(s, lg, store)
	go syncer.Start()

	b, err := balancer.New(ctx, appCFG, lg, store, tl, s, ses, sec, uc)
	if err != nil {
		lg.Fatal("failed to inizialise logic layer: %v", zap.String("error", err.Error()))
	}

	go runGRPC(ctx, lg, b, appCFG.Security, appCFG.Network, ses)

	// TODO: do check before connect
	time.Sleep(2 * time.Second)

	go runWebCLI(ctx, cfg.WebApp, cfg.Security, lg, cfg.Network.GRPC)

	if cfg.Network.REST != "" {
		go runREST(ctx, lg, b, cfg.Network)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cancel()
	time.Sleep(1 * time.Second)
}
