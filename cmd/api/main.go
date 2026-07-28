package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-clean-arch/internal/di"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config", "error", err)
	}
	server, err := di.InitializeAPI(cfg)
	if err != nil {
		log.Fatal("failed to initialize server", "error", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := server.Start(); err != nil {
			log.Error("HTTP server stopped unexpectedly", "error", err)
			stop()
		}
	}()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to gracefully shut down server", "error", err)
	}
}
