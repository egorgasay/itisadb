package main

import (
	"github.com/egorgasay/grpc-storage/config"
	"github.com/egorgasay/grpc-storage/internal/handler"
	"github.com/egorgasay/grpc-storage/internal/storage"
	"github.com/egorgasay/grpc-storage/internal/usecase"
	"github.com/egorgasay/grpc-storage/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/httplog"
)

func main() {
	cfg := config.New()

	store, err := storage.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	logic := usecase.New(store)
	handler := handler.New(logic, logger.New(log))
	router := chi.NewRouter()

	log := httplog.NewLogger("grpc-storage", httplog.Options{
		Concise: true,
	})
	router.Use(httplog.RequestLogger(log))
	router.Use(middleware.Recoverer)

	go func() {
		log.Info().Msg("Stating loyalty: " + cfg.Host)
		err := http.ListenAndServe(cfg.Host, router)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info().Msg("Shutdown Server ...")
}
