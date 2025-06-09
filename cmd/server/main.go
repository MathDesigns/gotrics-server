package main

import (
	"context"
	"errors"
	"gotrics-server/internal/api"
	"gotrics-server/internal/config"
	"gotrics-server/internal/logger"
	"gotrics-server/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appLogger := logger.New("SERVER: ")
	appLogger.Println("Starting Gotrics server...")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to load configuration: %v", err)
	}

	store := storage.NewInfluxStore(cfg.InfluxDB)
	defer store.Close()

	hub := api.NewHub()
	go hub.Run()

	apiServer := api.NewServer(appLogger, cfg, store, hub)

	srv := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: apiServer.Router(),
	}

	go func() {
		appLogger.Printf("Server listening on %s", cfg.ListenAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Fatalf("Could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Println("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	appLogger.Println("Server exiting.")
}
