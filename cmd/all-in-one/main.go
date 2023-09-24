package main

import (
	"context"
	"itisadb/cmd/uppers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	uppers.UpBalancer(ctx)
	uppers.UpGRPCStorage(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cancel()
	time.Sleep(2 * time.Second)
}
