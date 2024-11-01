package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/svirmi/binance-stream/storage"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading config, %v", err)
	}
}

func main() {

	logfile, err := os.Create("app.log")

	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(logfile, nil))

	dbClient, err := storage.NewMysqlClient(logger)

	defer func() {
		dbClient.Close()
		logfile.Close()
	}()

	if err != nil {
		logger.Error("error creating storage client", slog.String("storage", err.Error()))
		panic(err)
	}

	dbClient.CreateTableKLine()
}
