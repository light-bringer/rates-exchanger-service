package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/light-bringer/rates-exchanger-service/cron"
	"github.com/light-bringer/rates-exchanger-service/db"
	"github.com/light-bringer/rates-exchanger-service/internal/handler"
	"github.com/light-bringer/rates-exchanger-service/internal/service"
	"github.com/light-bringer/rates-exchanger-service/internal/sync"
)

func main() {
	dbParams := db.PostgresConfigParams{
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

	slog.Info("Starting the API service", "dbParams", dbParams)

	dbConfig := db.NewPostgresConfig(dbParams)

	if dbConfig == nil {
		log.Fatal("Invalid database configuration")
	}

	dbConn, err := db.BuildPGXConnPool(*dbConfig)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	defer dbConn.Close()
	syncService := sync.NewExchangeRateSync("rate_api", SyncURL, dbConn)

	go cron.PeriodicJob(syncService.Sync, SyncInterval)
	cleanSvc := func() {
		syncService.Cleanup(DeletionDays)
	}
	go cron.PeriodicJob(cleanSvc, DeleteInterval)
	ratesService := service.NewRatesService(dbConn, "rate_api")
	ratesHandler := handler.NewHandler(ratesService)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      ratesHandler.Routes(),
		ReadTimeout:  ServerTimeout,
		WriteTimeout: ServerTimeout,
	}

	log.Fatal(server.ListenAndServe())
}
