package main

import (
	"currency-exchange-converter/internal/config"
	"currency-exchange-converter/internal/db"
	"flag"
	"log"
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

	// 3. open sqlite
	database, err := db.Open(cfg.SQLite.Path)
	if err != nil {
		log.Fatal(err)
	}

	// 4. run migrations
	err = db.RunMigrations(database, "migrations/sqlite/000001_init.sql")
	if err != nil {
		log.Fatal(err)
	}

	// // 5. create logger
	// logger := logger.New(cfg.Log.Level)

	// // 6. start http server here (опущено)

	// // repository
	// exchangeRepo := repository.NewExchangeRateRepository(dbConn)

	// // service
	// exchangeService := service.NewExchangeRateService(exchangeRepo)

	// // handler
	// exchangeHandler := handler.NewExchangeRateHandler(exchangeService)

	// // router
	// router := handler.NewRouter(exchangeHandler)

	// // server
	// srv := server.New(cfg, logger, router)

	// log.Fatal(srv.Run())

}
