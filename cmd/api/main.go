package main

import (
	http "currency-exchange-converter/internal/api"
	"currency-exchange-converter/internal/config"
	"currency-exchange-converter/internal/currency"
	"currency-exchange-converter/internal/db"
	"currency-exchange-converter/internal/health"
	"currency-exchange-converter/internal/logger"
	"currency-exchange-converter/internal/server"
	"flag"
	"log"
	"os"
)

func main() {

	configPath := flag.String("config", "configs/config.local.yml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New(cfg.Log.Level)

	database, err := db.Open(cfg.SQLite.Path)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	err = db.RunMigrations(database, "migrations/sqlite/000001_init.sql")
	if err != nil {
		logger.Error("failed to run database migrations", "error", err)
		os.Exit(1)
	}

	currencyRepo := currency.NewRepository(database)
	currencyService := currency.NewService(currencyRepo)
	currencyHandler := currency.NewHandler(currencyService)

	healthService := health.NewService()
	healthHandler := health.NewHandler(healthService)

	router := http.NewRouter(healthHandler, currencyHandler)

	srv := server.New(cfg, logger, router)
	log.Fatal(srv.Run())

}
