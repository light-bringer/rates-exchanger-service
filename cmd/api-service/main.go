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

	"github.com/light-bringer/rates-exchanger-service/cron"
	"github.com/light-bringer/rates-exchanger-service/db"
	"github.com/light-bringer/rates-exchanger-service/internal/handler"
	"github.com/light-bringer/rates-exchanger-service/internal/service"
	"github.com/light-bringer/rates-exchanger-service/internal/sync"
	"github.com/light-bringer/rates-exchanger-service/models"
)

func main() {
	// read the configuration file
	configPath, paramErr := parseFlags()
	if paramErr != nil {
		log.Fatalf("Error parsing flags: %v", paramErr)
	}

	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Error reading the configuration file: %v", err)
	}

	// set the default values for the configuration
	setDefaults(config)

	dbParams := models.PostgresConfigParams{
		Host:           config.Database.Host,
		Port:           config.Database.Port,
		Username:       config.Database.User,
		Password:       config.Database.Pass,
		Database:       config.Database.Name,
		SSLMode:        string(config.Database.SSLMode),
		MinConnections: config.Database.MinConnections,
		MaxConnections: config.Database.MaxConnections,
		SchemaName:     config.Database.Schema,
	}

	// Create a context that listens for termination signals
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Kill,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGQUIT,
	)
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
	syncService := sync.NewExchangeRateSync(config.Database.Schema, config.CronJobs.Rates.SyncURL, dbConn)
	cleanSvc := func() {
		syncService.Cleanup(config.CronJobs.Cleanup.MaxAge)
	}
	// Run the cleanup service once before starting the cron job
	cleanSvc()

	// Run the sync service once before starting the cron job
	syncService.Sync()

	// Start the cron jobs
	go cron.Periodically(ctx, syncService.Sync, config.CronJobs.Rates.UpdateInterval)
	go cron.Periodically(ctx, cleanSvc, config.CronJobs.Cleanup.DeletionInterval)

	// Create a new rates service and handler
	ratesService := service.NewRatesService(dbConn, config.Database.Schema)
	ratesHandler := handler.NewHandler(ratesService)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.HTTP.Port),
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("Server exited gracefully!")
}
