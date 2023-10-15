package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ZeeeUs/BaumanGisSystem/internal/config"
	transport "github.com/ZeeeUs/BaumanGisSystem/internal/transport/http"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}
	logger := cfg.Logger()

	healthServer := transport.NewHealthServer(cfg.Server.HealthHost)
	go func() {
		if err = healthServer.Run(); err != nil {
			logger.Fatal().Err(err).Msg("health server starting error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := healthServer.Shutdown(context.Background()); err != nil {
			logger.Fatal().Err(err).Msg("health server shutdown error")
		}
		defer wg.Done()
	}()
}
