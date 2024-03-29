package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/light-bringer/rates-exchanger-service/cron"
	"github.com/light-bringer/rates-exchanger-service/db"
	"github.com/light-bringer/rates-exchanger-service/internal/handler"
	"github.com/light-bringer/rates-exchanger-service/internal/service"
	"github.com/light-bringer/rates-exchanger-service/internal/sync"
	"github.com/light-bringer/rates-exchanger-service/models"
)

func main() {
	dbParams := models.PostgresConfigParams{
		Host:           "localhost",
		Port:           5432,
		Username:       "postgres",
		Password:       "postgres",
		Database:       "postgres",
		SSLMode:        "disable",
		MinConnections: 1,
		MaxConnections: 10,
		SchemaName:     "public",
	}

	// Create a context that listens for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("Starting the API service", "dbParams", dbParams)

	dbConfig := db.NewPostgresConfig(dbParams)

	if dbConfig == nil {
		slog.Error("Error creating the database configuration")
		ctx.Done()
		return
	}

	dbConn, err := db.BuildPGXConnPool(context.Background(), *dbConfig)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	defer dbConn.Close()
	syncService := sync.NewExchangeRateSync("rate_api", SyncURL, dbConn)

	go cron.Periodically(ctx, syncService.Sync, SyncInterval)
	cleanSvc := func() {
		syncService.Cleanup(DeletionDays)
	}
	go cron.Periodically(ctx, cleanSvc, DeleteInterval)
	ratesService := service.NewRatesService(dbConn, "rate_api")
	ratesHandler := handler.NewHandler(ratesService)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      ratesHandler.Routes(),
		ReadTimeout:  ServerTimeout,
		WriteTimeout: ServerTimeout,
	}

	// Run server in a goroutine so that it doesn't block
	go func() {
		slog.Info(fmt.Sprintf("Starting server on %s\n", server.Addr))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Listen for the context cancellation signal from NotifyContext
	<-ctx.Done()

	// Context cancelled - proceed to gracefully shutdown the server
	slog.Info("Shutting down server...")
	// Perform any necessary cleanup after the context is canceled
	slog.Info("Received termination signal. Stopping cron job...")

	// Create a deadline to wait for.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("Server exited gracefully!")
}
