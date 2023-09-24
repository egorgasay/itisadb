package main

import (
	"golang.org/x/net/context"
	"itisadb/cmd/uppers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	uppers.UpGRPCStorage(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cancel()
	time.Sleep(1 * time.Second)
}
