package main

import (
	handler "currency-exchange-converter/internal/api"
	"currency-exchange-converter/internal/config"
	"currency-exchange-converter/internal/db"
	"currency-exchange-converter/internal/logger"
	"currency-exchange-converter/internal/server"
	"flag"
	"log"
	"os"
)

func main() {

	// 1. parse flags
	configPath := flag.String("config", "configs/config.local.yml", "path to config file")
	flag.Parse()

	// 2. load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// 3. create logger
	logger := logger.New(cfg.Log.Level)

	// 4. open sqlite
	database, err := db.Open(cfg.SQLite.Path)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("failed to close database", "error", err) // <-- Вот теперь можно так!
		}
	}()

	// 5. run migrations
	err = db.RunMigrations(database, "migrations/sqlite/000001_init.sql")
	if err != nil {
		logger.Error("failed to run database migrations", "error", err)
		os.Exit(1)
	}

	// repository
	//exchangeRepo := repository.NewExchangeRateRepository(dbConn)

	// service
	//exchangeService := service.NewExchangeRateService(exchangeRepo)

	// handler
	//exchangeHandler := handler.NewExchangeRateHandler(exchangeService)

	// router
	router := handler.NewRouter()

	// 6. start http server here (опущено)
	srv := server.New(cfg, logger, router)

	log.Fatal(srv.Run())

}
